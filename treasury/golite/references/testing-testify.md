# Testing: table-driven testify + mockery

The default Go backend test pattern. Conform to the repo if it already tests differently (gomock, stdlib); otherwise use this.

## The pattern, in pieces
- **External test package** `package x_test` — tests the package as a client, through its exported surface.
- **`r := require.New(t)`**, created inside each `t.Run` subtest so it binds that subtest's `t` — `require` fails fast (stops the case on first failed assertion). Use `assert` only when you want several soft checks. `assert.AnError` is a ready-made generic error for failure cases.
- **In-function structs**: `mockArgs` holds the mocks; `args` holds the call inputs (`ctx`, request/message); `want` (or a plain `wantErr bool`) holds the expected outcome.
- **`prepare func(m mockArgs, a args)`** — per-case closure that programs mock expectations. An empty closure means "no calls expected."
- **Case naming** — `"success, case ..."` and `"fail, case ..."`. One case per behavior branch: the happy path plus each failure mode (bad input, dependency error, downstream error).
- **Loop** — `for _, tt := range tests { t.Run(tt.name, func(t *testing.T) { ... }) }`.

## Mocks: generated, never hand-written
- mockery with the testify template, configured in `.mockery.yaml`. Mocks land in `<pkg>/mocks/` as generated `mocks.go` (`package <pkg>_mocks`, struct `<Iface>Mock`).
- Constructor `NewXMock(t)` registers a cleanup that asserts all expectations were met — no manual `AssertExpectations`.
- Expecter API: `m.X.EXPECT().Method(args).Return(vals)`. Match arguments with `mock.MatchedBy(func(v T) bool { ... })`, or `mock.Anything` when the argument is irrelevant.
- When you change an interface, **regenerate** (`mockery`) — never edit generated `mocks.go`.

## Dependency injection
Build the unit under test with its constructor config, passing mocks: `NewHandler(HandlerConfig{ProductStorage: m.productStorage})`. This is why access is an interface.

## HTTP handler variant
`httptest.NewRecorder()` + the framework's test context (e.g. `gin.CreateTestContext(w)`), build the request, call the handler, decode the response envelope, assert status + code + message + data. `// Arrange / Act / Assert` markers are part of this convention.

## Consumer / non-HTTP variant
Call the method directly — `err := h.OnProductCreated(ctx, msg)` — and assert `r.Error(err)` / `r.NoError(err)`.

## Coverage and command
One case per branch. Run:
```
go test -v -race -covermode=atomic -coverpkg=./... ./...
```

## Worked example (handler)
```go
package product_test

func TestCreateProduct(t *testing.T) {
	type mockArgs struct{ store *access_mocks.ProductStorageMock }
	type args struct {
		ctx context.Context
		req product.CreateRequest
	}
	type want struct {
		status  int
		errCode app.Code
	}

	tests := []struct {
		name    string
		prepare func(m mockArgs, a args)
		args    args
		want    want
	}{
		{
			name: "success, case valid request",
			prepare: func(m mockArgs, a args) {
				m.store.EXPECT().
					Create(a.ctx, mock.MatchedBy(func(p access.Product) bool { return p.Name == a.req.Name })).
					Return(access.Product{ID: "p-1"}, nil)
			},
			args: args{req: product.CreateRequest{Name: "Widget"}},
			want: want{status: http.StatusCreated, errCode: app.CodeSuccess},
		},
		{
			name:    "fail, case missing name",
			prepare: func(m mockArgs, a args) {}, // no calls
			args:    args{req: product.CreateRequest{}},
			want:    want{status: http.StatusBadRequest, errCode: app.CodeBadRequest},
		},
		{
			name: "fail, case storage error",
			prepare: func(m mockArgs, a args) {
				m.store.EXPECT().Create(a.ctx, mock.Anything).Return(access.Product{}, assert.AnError)
			},
			args: args{req: product.CreateRequest{Name: "Widget"}},
			want: want{status: http.StatusInternalServerError, errCode: app.CodeInternalError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			// Arrange
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tt.args.req)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/product", &body)
			tt.args.ctx = ctx.Request.Context()

			m := mockArgs{store: access_mocks.NewProductStorageMock(t)}
			tt.prepare(m, tt.args)

			h := product.NewHandler(product.HandlerConfig{ProductStorage: m.store})

			// Act
			h.CreateProduct(ctx)

			// Assert
			var resp wrapper.ResponseOption[product.CreateResponse]
			json.NewDecoder(w.Body).Decode(&resp)
			r.Equal(tt.want.status, w.Code)
			r.Equal(tt.want.errCode, resp.Code)
		})
	}
}
```
Repo-specific bits above (`gin`, `wrapper.ResponseOption[T]`, `app.Code*`, `access_mocks`) are illustrative — swap them for whatever the target repo uses. The portable core is: external `_test` package, `require`, the `mockArgs`/`args`/`want` table, the `prepare` closure, mockery expecter mocks, one case per branch.
