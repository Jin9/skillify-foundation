# Test Patterns

Every `app/` package targets **100% statement coverage**. Every `if err != nil` branch — including model-getter parse failures and `kafka.BindMessage` validation failures — needs a dedicated case.

---

## 1. HTTP handler test

Canonical pattern with `mockArgs` / `args` / `want` local types and a `prepare` function.

```go
func TestCreateProduct(t *testing.T) {
    r := require.New(t)

    productID := uuid.New()

    type mockArgs struct {
        productStorage *access_mocks.ProductStorageMock
    }

    type args struct {
        ctx context.Context
        req product.CreateProductRequest
    }

    type want struct {
        err     bool
        code    wrapper.Code
        Message wrapper.Message
    }

    tests := []struct {
        name    string
        prepare func(m mockArgs, args args)
        args    args
        want    want
    }{
        {
            name: "success, case valid request",
            prepare: func(m mockArgs, args args) {
                m.productStorage.
                    EXPECT().
                    CreateProduct(args.ctx, mock.MatchedBy(func(p access.Product) bool {
                        return p.Name == args.req.Name && p.Price == args.req.Price
                    })).
                    Return(access.Product{ProductID: productID.String()}, nil)
            },
            args: args{
                req: product.CreateProductRequest{
                    Name:           "Test Product",
                    Price:          99.99,
                    OrganizationID: uuid.New().String(),
                },
            },
            want: want{
                err:     false,
                code:    wrapper.CodeSuccess,
                Message: wrapper.MessageSuccess,
            },
        },
        {
            name:    "fail, case invalid request body - missing name",
            prepare: func(m mockArgs, args args) {},
            args:    args{req: product.CreateProductRequest{Price: 99.99}},
            want: want{
                err:     true,
                code:    wrapper.CodeBadRequest,
                Message: wrapper.MessageBadRequest,
            },
        },
        // ... one case per error branch for 100% coverage
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            ctx, _ := gin.CreateTestContext(w)

            var payload bytes.Buffer
            json.NewEncoder(&payload).Encode(tt.args.req)

            req := httptest.NewRequest(http.MethodPost, "http://0.0.0.0/api/v1/product", &payload)
            req.Header.Set("Content-Type", "application/json")
            ctx.Request = req

            m := mockArgs{
                productStorage: access_mocks.NewProductStorageMock(t),
            }

            if tt.prepare != nil {
                tt.args.ctx = ctx.Request.Context() // assign before prepare
                tt.prepare(m, tt.args)
            }

            h := product.NewHandler(product.HandlerConfig{
                ProductStorage: m.productStorage,
            })

            h.CreateProduct(ctx)

            var resp wrapper.ResponseOption[product.CreateProductResponse]
            json.NewDecoder(w.Body).Decode(&resp)

            if tt.want.err {
                r.NotEqual(http.StatusOK, w.Code)
                r.Equal(tt.want.code, resp.Code)
                r.Equal(tt.want.Message, resp.Message)
            } else {
                r.Equal(http.StatusCreated, w.Code)
                r.Equal(tt.want.code, resp.Code)
                r.Equal(tt.want.Message, resp.Message)
                r.NotNil(resp.Data)
                r.Equal(productID, resp.Data.ProductID)
            }
        })
    }
}
```

Key points:

- **Local type declarations** inside the test function: `mockArgs`, `args`, `want`.
- **`prepare` function** (not `setup`) receives `mockArgs` and `args` — sets up mock expectations using data from `args`.
- `args.ctx` is assigned inside the loop from `ctx.Request.Context()` **before** calling `prepare`.
- Use `gin.CreateTestContext(httptest.NewRecorder())` — returns `(*gin.Context, *gin.Engine)`.
- Decode response into `wrapper.ResponseOption[T]`.
- Assert with `require.New(t)` shorthand `r`.
- Test every `if err != nil` branch, including model-getter failures (e.g. `GetID()` with `"not-a-valid-uuid"`).

---

## 2. Kafka consumer test

```go
func TestOnProductCreated(t *testing.T) {
    r := require.New(t)

    type mockArgs struct {
        productStorage *access_mocks.ProductStorageMock
    }

    type args struct {
        ctx context.Context
        msg kafka.Message[json.RawMessage]
    }

    tests := []struct {
        name    string
        prepare func(m mockArgs, args args)
        args    args
        wantErr bool
    }{
        {
            name: "success, case valid event",
            prepare: func(m mockArgs, args args) {
                m.productStorage.EXPECT().
                    CreateProduct(args.ctx, mock.Anything).
                    Return(access.Product{}, nil)
            },
            args: args{
                ctx: context.Background(),
                msg: kafka.Message[json.RawMessage]{
                    EventName:   "PRODUCT_CREATED",
                    AggregateID: uuid.New().String(),
                    EventID:     uuid.New().String(),
                    Timestamp:   time.Now(),
                    Payload: func() json.RawMessage {
                        b, _ := json.Marshal(product.CreateProductMessage{
                            Name:           "Test Product",
                            Price:          10.00,
                            OrganizationID: "org-123",
                        })
                        return b
                    }(),
                },
            },
            wantErr: false,
        },
        {
            name:    "fail, case invalid JSON payload",
            prepare: func(m mockArgs, args args) {},
            args: args{
                ctx: context.Background(),
                msg: kafka.Message[json.RawMessage]{
                    Payload: json.RawMessage(`{invalid`),
                },
            },
            wantErr: true,
        },
        {
            name:    "fail, case validation error - missing required fields",
            prepare: func(m mockArgs, args args) {},
            args: args{
                ctx: context.Background(),
                msg: kafka.Message[json.RawMessage]{
                    Payload: json.RawMessage(`{"description":"no name or price"}`),
                },
            },
            wantErr: true,
        },
        {
            name: "fail, case storage error",
            prepare: func(m mockArgs, args args) {
                m.productStorage.EXPECT().
                    CreateProduct(args.ctx, mock.Anything).
                    Return(access.Product{}, assert.AnError)
            },
            args: args{
                ctx: context.Background(),
                msg: kafka.Message[json.RawMessage]{
                    Payload: func() json.RawMessage {
                        b, _ := json.Marshal(product.CreateProductMessage{
                            Name: "Test", Price: 10.00, OrganizationID: "org-123",
                        })
                        return b
                    }(),
                },
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := mockArgs{
                productStorage: access_mocks.NewProductStorageMock(t),
            }
            if tt.prepare != nil {
                tt.prepare(m, tt.args)
            }

            h := product.NewHandler(product.HandlerConfig{
                ProductStorage: m.productStorage,
            })

            err := h.OnProductCreated(tt.args.ctx, tt.args.msg)
            if tt.wantErr {
                r.Error(err)
            } else {
                r.NoError(err)
            }
        })
    }
}
```

Required cases for every consumer test:

1. Success path.
2. Invalid JSON payload (e.g. `json.RawMessage("{invalid")`).
3. Validation failure — payload missing required fields enforced by binding tags.
4. Storage / downstream error.

---

## 3. Coverage strategy

- Run: `make test` → `go test -v -race -covermode=atomic -buildvcs -coverpkg=./... ./...`.
- Every `if err != nil` branch must have a dedicated test case.
- Use `mock.MatchedBy(func(...) bool)` for complex argument matching.
- Model-getter errors (`GetID()` with `"not-a-valid-uuid"`) need test cases.
- `kafka.BindMessage` validation failures need test cases (missing required fields).

---

## 4. Mock generation

Mocks are configured in `.mockery.yaml` at the repo root:

```yaml
all: false
recursive: true
dir: "{{.InterfaceDir}}/mocks"
filename: "mocks.go"
force-file-write: true
formatter: goimports
structname: "{{.InterfaceName}}Mock"
pkgname: "{{ .SrcPackageName }}_mocks"
template: testify
packages:
  gitdev.devops.krungthai.com/ktb-digital-lending/dgl-new/shareable/go-template:
    config:
      all: true
```

| Aspect | Value |
|---|---|
| Mock struct naming | `<InterfaceName>Mock` (e.g. `ProductStorageMock`) |
| Package naming | `<srcpkg>_mocks` (e.g. `access_mocks`) |
| Output | `<interface-dir>/mocks/mocks.go` |
| Constructor | `access_mocks.NewProductStorageMock(t)` — auto-registers cleanup |
| Fluent API | `mock.EXPECT().Method(args...).Return(vals...)` |

Regenerate all mocks:

```bash
mockery
```
