# Testing pattern

Every business-logic file ships with a matching `_test.go`. Coverage target is **100% statement coverage** in `app/<domain>/` (excluding generated `mocks/`).

## The shape: `mockArgs` / `args` / `want` / `prepare` + table

The pattern is identical for handlers and consumers — only the act/assert blocks differ.

```go
package <domain>_test

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "gitlab.com/.../common/app"
    "gitlab.com/.../common/wrapper"
    access_mocks "gitlab.com/.../app/<domain>/access/mocks"

    "gitlab.com/.../app/<domain>"
    "gitlab.com/.../app/<domain>/access"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

func Test<Action>(t *testing.T) {
    r := require.New(t)

    // Local types — keep them inside the test function for scope clarity.
    type mockArgs struct {
        <dep>Storage *access_mocks.<Dep>StorageMock
        // add other mocks as needed
    }

    type args struct {
        ctx context.Context
        req <domain>.<Action>Request
    }

    type want struct {
        err     bool
        code    app.Code
        Message app.Message
        data    wrapper.ResponseOption[<domain>.<Action>Response]
    }

    tests := []struct {
        name    string
        prepare func(m mockArgs, args args)
        args    args
        want    want
    }{
        // success
        {
            name: "success, case valid request",
            prepare: func(m mockArgs, args args) {
                m.<dep>Storage.
                    EXPECT().
                    <Method>(args.ctx, mock.Anything).
                    Return(access.<Model>{<ModelID>: uuid.New().String()}, nil)
            },
            args: args{
                req: <domain>.<Action>Request{ /* valid fields */ },
            },
            want: want{
                err:     false,
                code:    app.CodeSuccess,
                Message: app.MessageSuccess,
            },
        },

        // each missing-required-field case
        {
            name: "fail, case invalid request body - missing <field>",
            prepare: func(m mockArgs, args args) { /* no storage call expected */ },
            args: args{
                req: <domain>.<Action>Request{ /* missing the field */ },
            },
            want: want{
                err:     true,
                code:    app.CodeBadRequest,
                Message: app.MessageBadRequest,
            },
        },

        // each downstream error
        {
            name: "fail, case <Method> storage error",
            prepare: func(m mockArgs, args args) {
                m.<dep>Storage.
                    EXPECT().
                    <Method>(args.ctx, mock.Anything).
                    Return(access.<Model>{}, assert.AnError)
            },
            args: args{ req: <domain>.<Action>Request{ /* valid */ } },
            want: want{
                err:     true,
                code:    app.CodeInternalError,
                Message: app.MessageInternalError,
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            ctx, _ := gin.CreateTestContext(w)

            // Arrange request body
            var payload bytes.Buffer
            json.NewEncoder(&payload).Encode(tt.args.req)
            req := httptest.NewRequest(http.MethodPost, "http://0.0.0.0/api/v1/platform/<domain>/<action>", &payload)
            req.Header.Set("Content-Type", "application/json")
            ctx.Request = req

            // Construct mocks
            m := mockArgs{
                <dep>Storage: access_mocks.New<Dep>StorageMock(t),
            }

            // IMPORTANT: assign ctx BEFORE calling prepare — mock matchers may capture it.
            if tt.prepare != nil {
                tt.args.ctx = ctx.Request.Context()
                tt.prepare(m, tt.args)
            }

            h := <domain>.NewHandler(<domain>.HandlerConfig{
                <Dep>Storage: m.<dep>Storage,
            })

            // Act
            h.<Action>(ctx)

            // Assert
            var resp wrapper.ResponseOption[<domain>.<Action>Response]
            json.NewDecoder(w.Body).Decode(&resp)

            if tt.want.err {
                r.NotEqual(http.StatusOK, w.Code)
                r.Equal(tt.want.code, resp.Code)
                r.Equal(tt.want.Message, resp.Message)
            } else {
                r.Equal(http.StatusOK, w.Code)
                r.Equal(tt.want.code, resp.Code)
                r.Equal(tt.want.Message, resp.Message)
            }
        })
    }
}
```

## Kafka consumer variant

Replace the Act/Assert section with:

```go
// Arrange — marshal payload into a kafka.Message[json.RawMessage]
payloadBytes, _ := json.Marshal(tt.args.payload)
msg := kafka.Message[json.RawMessage]{
    EventID: "test-event-id",
    Payload: payloadBytes,
}

// Construct handler
h := <domain>.NewHandler(<domain>.HandlerConfig{
    <Dep>Storage: m.<dep>Storage,
})

// Act
err := h.On<Action>(tt.args.ctx, msg)

// Assert
if tt.want.err {
    r.Error(err)
} else {
    r.NoError(err)
}
```

For invalid-JSON cases, set `msg.Payload = []byte("{ invalid")` and expect an error.
For validation-failure cases, marshal a struct that omits a `binding:"required"` field and expect an error.

## Coverage rules

These branches MUST have a test case each:

- Every `if err != nil` block in the handler/consumer/service body.
- Every `binding:"required"` field on the request/payload struct (one missing-field case per).
- Every model-getter that returns `(T, error)` — e.g. `GetID()` parsing `MemberID` as UUID. Test with a deliberately bad string.
- For Kafka consumers: invalid JSON (`{ invalid`) and validation-failure cases.
- The success path.

Verify with:

```bash
go test -race -coverprofile=coverage.out ./app/<domain>/...
go tool cover -func=coverage.out | grep -v 100.0%
```

The grep should print nothing (every function at 100%).

## Mock-builder API

Mocks are generated by `mockery` per the repo's `.mockery.yaml`. **Do not hand-edit** them.

Pattern:

```go
m := access_mocks.NewMemberStorageMock(t)

m.EXPECT().
    GetMemberByEmail(ctx, hashedEmail).
    Return(access.Member{MemberID: id.String()}, nil)

m.EXPECT().
    GetMemberByEmail(mock.Anything, mock.MatchedBy(func(s string) bool {
        return len(s) > 0
    })).
    Return(access.Member{}, assert.AnError)
```

After changing an interface, regenerate mocks:

```bash
mockery
```

(or `go generate ./...` if a `//go:generate` directive is in use).

## Common pitfalls

| Symptom | Cause | Fix |
|---|---|---|
| Test stalls on `mock.MatchedBy` | `args.ctx` not assigned before `prepare` is called | Move `tt.args.ctx = ctx.Request.Context()` BEFORE `tt.prepare(m, tt.args)`. |
| Mock returns wrong type | Stale mock after interface change | Regenerate with `mockery`. |
| Coverage stuck at 95% | Missing test for a `GetID()` parse failure | Add a case where the storage returns a model with a deliberately invalid UUID string. |
| `json.Unmarshal` worked in handler test but consumer rejects valid payload | Consumer uses `kafka.BindMessage` (correct) which runs binding tags | Verify the test marshals a payload that satisfies all `binding` tags. |
| `make test` passes but `make precommit` fails on lint | Mocks file unformatted, or test file imports out of order | Run `go fmt` and `goimports`. Do not modify `.golangci.yaml` to silence. |
