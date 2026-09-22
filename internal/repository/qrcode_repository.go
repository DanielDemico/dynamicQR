package repository

import (
	"context"
	"errors"
	"time"

	"dynamicqr/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrNotFound     = errors.New("QR Code não encontrado")
	ErrSlugConflict = errors.New("O slug informado já está em uso")
)

type QRCodeRepository struct {
	collection *mongo.Collection
}

func NewQRCodeRepository(db *mongo.Database) *QRCodeRepository {
	return &QRCodeRepository{
		collection: db.Collection("qr_codes"),
	}
}

// Create insere um novo QR Code garantindo timestamps e valores iniciais
func (r *QRCodeRepository) Create(ctx context.Context, qr *domain.QRCode) error {
	now := time.Now().UTC()
	if qr.CreatedAt.IsZero() {
		qr.CreatedAt = now
	}
	qr.UpdatedAt = now
	qr.IsActive = true
	qr.ClickCount = 0
	if qr.QRMimeType == "" {
		qr.QRMimeType = "image/png"
	}

	result, err := r.collection.InsertOne(ctx, qr)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrSlugConflict
		}
		return err
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		qr.ID = oid
	}

	return nil
}

// FindByID busca um QR Code pelo seu ObjectID hexadecimal
func (r *QRCodeRepository) FindByID(ctx context.Context, idHex string) (*domain.QRCode, error) {
	oid, err := bson.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, ErrNotFound
	}

	var qr domain.QRCode
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&qr)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &qr, nil
}

// FindBySlug busca um QR Code ativo pelo slug exclusivo
func (r *QRCodeRepository) FindBySlug(ctx context.Context, slug string) (*domain.QRCode, error) {
	var qr domain.QRCode
	err := r.collection.FindOne(ctx, bson.M{"slug": slug, "is_active": true}).Decode(&qr)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &qr, nil
}

// List retorna QR Codes ordenados por data decrescente com paginação e busca
func (r *QRCodeRepository) List(ctx context.Context, search string, page, limit int64) ([]domain.QRCode, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	filter := bson.M{}
	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"title": bson.M{"$regex": search, "$options": "i"}},
				{"slug": bson.M{"$regex": search, "$options": "i"}},
				{"target_url": bson.M{"$regex": search, "$options": "i"}},
			},
		}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * limit
	findOpts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var list []domain.QRCode
	if err := cursor.All(ctx, &list); err != nil {
		return nil, 0, err
	}

	if list == nil {
		list = []domain.QRCode{}
	}

	return list, total, nil
}

// Update atualiza dinamicamente título, target_url ou is_active preservando slug e imagem
func (r *QRCodeRepository) Update(ctx context.Context, idHex string, title, targetURL *string, isActive *bool) (*domain.QRCode, error) {
	oid, err := bson.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, ErrNotFound
	}

	updateFields := bson.M{
		"updated_at": time.Now().UTC(),
	}
	if title != nil {
		updateFields["title"] = *title
	}
	if targetURL != nil {
		updateFields["target_url"] = *targetURL
	}
	if isActive != nil {
		updateFields["is_active"] = *isActive
	}

	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated domain.QRCode
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": updateFields}, opt).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &updated, nil
}

// Delete remove o documento por ID
func (r *QRCodeRepository) Delete(ctx context.Context, idHex string) error {
	oid, err := bson.ObjectIDFromHex(idHex)
	if err != nil {
		return ErrNotFound
	}

	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// IncrementClickCount incrementa atomicamente o contador de acessos do slug
func (r *QRCodeRepository) IncrementClickCount(ctx context.Context, slug string) error {
	update := bson.M{
		"$inc": bson.M{"click_count": 1},
		"$set": bson.M{"updated_at": time.Now().UTC()},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"slug": slug}, update)
	return err
}
