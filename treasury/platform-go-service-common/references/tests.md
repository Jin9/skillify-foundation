# Test Patterns

Table-driven tests with mockery `.EXPECT()` mocks for HTTP handlers and Kafka consumers. Target 100% statement coverage in every `app/` package.

## 1. HTTP Handler Test

```go
func TestHandler_CreateProduct(t *testing.T) {
    tests := []struct {
        name     string
        body     any
        setup    func(m *access_mocks.StorageMock, p *kafka_mocks.ProducerMock)
        wantCode int
        wantBody wrapper.Response[*ProductResponse]
    }{
        {
            name: "success",
            body: CreateProductRequest{Name: "Widget"},
            setup: func(m *access_mocks.StorageMock, p *kafka_mocks.ProducerMock) {
                m.EXPECT().Insert(mock.Anything, mock.Anything).
                    Return(&access.Model{ID: "1", Name: "Widget"}, nil)
                p.EXPECT().SendMessage(mock.Anything, mock.Anything, mock.Anything).
                    Return(nil)
            },
            wantCode: http.StatusCreated,
            wantBody: wrapper.Response[*ProductResponse]{
                Code:    app.CodeProductCreated,
                Message: app.MessageProductCreated,
                Data:    &ProductResponse{ID: "1", Name: "Widget"},
            },
        },
        {
            name:     "invalid body",
            body:     "invalid",
            setup:    func(m *access_mocks.StorageMock, p *kafka_mocks.ProducerMock) {},
            wantCode: http.StatusBadRequest,
        },
        // ... one test case per error branch for 100% coverage
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockStorage := access_mocks.NewStorageMock(t)
            mockProducer := kafka_mocks.NewProducerMock(t)
            tt.setup(mockStorage, mockProducer)

            h := NewHandler(HandlerConfig{
                Storage:  mockStorage,
                Producer: mockProducer,
                Clock:    clock.New(),
            })

            bodyBytes, _ := json.Marshal(tt.body)
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)
            c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bodyBytes))
            c.Request.Header.Set("Content-Type", "application/json")

            h.CreateProduct(c)

            require.Equal(t, tt.wantCode, w.Code)
            if tt.wantBody.Code != "" {
                var got wrapper.Response[*ProductResponse]
                require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
                assert.Equal(t, tt.wantBody, got)
            }
        })
    }
}
```

## 2. Kafka Consumer Test

```go
func TestHandler_OnProductCreated(t *testing.T) {
    tests := []struct {
        name    string
        msg     kafka.Message[json.RawMessage]
        setup   func(m *access_mocks.StorageMock)
        wantErr bool
    }{
        {
            name: "success",
            msg:  kafka.Message[json.RawMessage]{Payload: json.RawMessage(`{"id":"1","name":"Widget"}`)},
            setup: func(m *access_mocks.StorageMock) {
                m.EXPECT().Upsert(mock.Anything, mock.Anything).Return(nil)
            },
        },
        {
            name:    "unmarshal error",
            msg:     kafka.Message[json.RawMessage]{Payload: json.RawMessage(`invalid`)},
            setup:   func(m *access_mocks.StorageMock) {},
            wantErr: true,
        },
        {
            name: "storage error",
            msg:  kafka.Message[json.RawMessage]{Payload: json.RawMessage(`{"id":"1"}`)},
            setup: func(m *access_mocks.StorageMock) {
                m.EXPECT().Upsert(mock.Anything, mock.Anything).Return(assert.AnError)
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := access_mocks.NewStorageMock(t)
            tt.setup(m)
            h := NewHandler(HandlerConfig{Storage: m})
            err := h.OnProductCreated(context.Background(), tt.msg)
            if tt.wantErr {
                require.Error(t, err)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

## 3. Coverage Rules

- Target 100% statement coverage for all `app/` packages.
- Every `if err != nil` branch must have a dedicated test case.
- Use `mock.MatchedBy(func(...) bool)` for complex argument matching.
- Model getter errors need test cases with invalid inputs.
