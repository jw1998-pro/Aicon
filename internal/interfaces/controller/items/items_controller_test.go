package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"Aicon-assignment/internal/domain/entity"
	domainErrors "Aicon-assignment/internal/domain/errors"
	"Aicon-assignment/internal/usecase"
)

// MockItemUsecase はtestify/mockを使用したモックユースケース
type MockItemUsecase struct {
	mock.Mock
}

func (m *MockItemUsecase) GetAllItems(ctx context.Context) ([]*entity.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Item), args.Error(1)
}

func (m *MockItemUsecase) GetItemByID(ctx context.Context, id int64) (*entity.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemUsecase) CreateItem(ctx context.Context, input usecase.CreateItemInput) (*entity.Item, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemUsecase) UpdateItem(ctx context.Context, id int64, input usecase.UpdateItemInput) (*entity.Item, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemUsecase) DeleteItem(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemUsecase) GetCategorySummary(ctx context.Context) (*usecase.CategorySummary, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.CategorySummary), args.Error(1)
}

func TestItemHandler_UpdateItem(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		requestBody    map[string]interface{}
		setupMock      func(*MockItemUsecase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "正常系: 名前のみ更新",
			id:   "1",
			requestBody: map[string]interface{}{
				"name": "更新された名前",
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				item, _ := entity.NewItem("更新された名前", "時計", "ROLEX", 1000000, "2023-01-01")
				item.ID = 1
				mockUsecase.On("UpdateItem", mock.Anything, int64(1), usecase.UpdateItemInput{
					Name: stringPtr("更新された名前"),
				}).Return(item, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":             float64(1),
				"name":           "更新された名前",
				"category":       "時計",
				"brand":          "ROLEX",
				"purchase_price": float64(1000000),
				"purchase_date":  "2023-01-01",
			},
		},
		{
			name: "正常系: ブランドのみ更新",
			id:   "1",
			requestBody: map[string]interface{}{
				"brand": "更新されたブランド",
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				item, _ := entity.NewItem("時計1", "時計", "更新されたブランド", 1000000, "2023-01-01")
				item.ID = 1
				mockUsecase.On("UpdateItem", mock.Anything, int64(1), usecase.UpdateItemInput{
					Brand: stringPtr("更新されたブランド"),
				}).Return(item, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":             float64(1),
				"name":           "時計1",
				"category":       "時計",
				"brand":          "更新されたブランド",
				"purchase_price": float64(1000000),
				"purchase_date":  "2023-01-01",
			},
		},
		{
			name: "正常系: 購入価格のみ更新",
			id:   "1",
			requestBody: map[string]interface{}{
				"purchase_price": 2000000,
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				item, _ := entity.NewItem("時計1", "時計", "ROLEX", 2000000, "2023-01-01")
				item.ID = 1
				mockUsecase.On("UpdateItem", mock.Anything, int64(1), usecase.UpdateItemInput{
					PurchasePrice: intPtr(2000000),
				}).Return(item, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":             float64(1),
				"name":           "時計1",
				"category":       "時計",
				"brand":          "ROLEX",
				"purchase_price": float64(2000000),
				"purchase_date":  "2023-01-01",
			},
		},
		{
			name: "正常系: 複数フィールド更新",
			id:   "1",
			requestBody: map[string]interface{}{
				"name":           "更新された名前",
				"brand":          "更新されたブランド",
				"purchase_price": 2000000,
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				item, _ := entity.NewItem("更新された名前", "時計", "更新されたブランド", 2000000, "2023-01-01")
				item.ID = 1
				mockUsecase.On("UpdateItem", mock.Anything, int64(1), usecase.UpdateItemInput{
					Name:          stringPtr("更新された名前"),
					Brand:         stringPtr("更新されたブランド"),
					PurchasePrice: intPtr(2000000),
				}).Return(item, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":             float64(1),
				"name":           "更新された名前",
				"category":       "時計",
				"brand":          "更新されたブランド",
				"purchase_price": float64(2000000),
				"purchase_date":  "2023-01-01",
			},
		},
		{
			name: "異常系: 無効なID",
			id:   "invalid",
			requestBody: map[string]interface{}{
				"name": "更新された名前",
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				// UpdateItemは呼ばれない
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid item ID",
			},
		},
		{
			name: "異常系: 存在しないアイテム",
			id:   "999",
			requestBody: map[string]interface{}{
				"name": "更新された名前",
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				mockUsecase.On("UpdateItem", mock.Anything, int64(999), usecase.UpdateItemInput{
					Name: stringPtr("更新された名前"),
				}).Return((*entity.Item)(nil), domainErrors.ErrItemNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "item not found",
			},
		},
		{
			name: "異常系: 空の名前",
			id:   "1",
			requestBody: map[string]interface{}{
				"name": "",
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				// UpdateItemは呼ばれない（バリデーションでエラー）
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "validation failed",
				"details": []interface{}{
					"name cannot be empty",
				},
			},
		},
		{
			name: "異常系: 負の購入価格",
			id:   "1",
			requestBody: map[string]interface{}{
				"purchase_price": -1,
			},
			setupMock: func(mockUsecase *MockItemUsecase) {
				// UpdateItemは呼ばれない（バリデーションでエラー）
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "validation failed",
				"details": []interface{}{
					"purchase_price must be 0 or greater",
				},
			},
		},
		{
			name: "異常系: フィールドが提供されていない",
			id:   "1",
			requestBody: map[string]interface{}{},
			setupMock: func(mockUsecase *MockItemUsecase) {
				// UpdateItemは呼ばれない（バリデーションでエラー）
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "validation failed",
				"details": []interface{}{
					"at least one field (name, brand, or purchase_price) must be provided",
				},
			},
		},
		{
			name: "異常系: 無効なJSON",
			id:   "1",
			requestBody: nil, // 無効なJSONを送信
			setupMock: func(mockUsecase *MockItemUsecase) {
				// UpdateItemは呼ばれない
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid request format",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Echoインスタンスを作成
			e := echo.New()
			
			// モックユースケースを作成
			mockUsecase := new(MockItemUsecase)
			tt.setupMock(mockUsecase)
			
			// ハンドラーを作成
			handler := NewItemHandler(mockUsecase)

			// リクエストボディを作成
			var reqBody []byte
			var err error
			if tt.requestBody != nil {
				reqBody, err = json.Marshal(tt.requestBody)
				require.NoError(t, err)
			} else {
				reqBody = []byte("invalid json")
			}

			// リクエストを作成
			req := httptest.NewRequest(http.MethodPatch, "/items/"+tt.id, bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			
			// レスポンスレコーダーを作成
			rec := httptest.NewRecorder()
			
			// Echoコンテキストを作成
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			// ハンドラーを実行
			err = handler.UpdateItem(c)

			// レスポンスを検証
			if tt.expectedStatus == http.StatusOK {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.expectedStatus, rec.Code)

			// レスポンスボディを検証
			var responseBody map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			// 期待されるフィールドを検証
			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, responseBody[key])
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

// ヘルパー関数
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
} 