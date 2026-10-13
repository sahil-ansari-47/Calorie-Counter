package auth

import (
	"caloriecounter/api/internal/database/sqlc"
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrInvalidCredentials = errors.New("invalid credentials")

type Profile struct {
	ID                               uuid.UUID  `json:"id"`
	Email                            string     `json:"email"`
	DisplayName                      *string    `json:"display_name,omitempty"`
	DateOfBirth                      *time.Time `json:"date_of_birth,omitempty"`
	Age                              *int32     `json:"age,omitempty"`
	HeightCm                         *float64   `json:"height_cm,omitempty"`
	Sex                              *string    `json:"sex,omitempty"`
	ActivityLevel                    *string    `json:"activity_level,omitempty"`
	Timezone                         string     `json:"timezone"`
	MaintenanceCalories              *float64   `json:"maintenance_calories,omitempty"`
	RecommendedCalories              *float64   `json:"recommended_calories,omitempty"`
	RecommendedProteinGDay           *float64   `json:"recommended_protein_g_day,omitempty"`
	RecommendedCarbohydratesGDay     *float64   `json:"recommended_carbohydrates_g_day,omitempty"`
	RecommendedFatGDay               *float64   `json:"recommended_fat_g_day,omitempty"`
	RecommendedFiberGDay             *float64   `json:"recommended_fiber_g_day,omitempty"`
	RecommendedSugarsGDay            *float64   `json:"recommended_sugars_g_day,omitempty"`
	RecommendedSaturatedFatGDay      *float64   `json:"recommended_saturated_fat_g_day,omitempty"`
	RecommendedCalciumMgWeek         *float64   `json:"recommended_calcium_mg_week,omitempty"`
	RecommendedIronMgWeek            *float64   `json:"recommended_iron_mg_week,omitempty"`
	RecommendedMagnesiumMgWeek       *float64   `json:"recommended_magnesium_mg_week,omitempty"`
	RecommendedPhosphorusMgWeek      *float64   `json:"recommended_phosphorus_mg_week,omitempty"`
	RecommendedPotassiumMgWeek       *float64   `json:"recommended_potassium_mg_week,omitempty"`
	RecommendedZincMgWeek            *float64   `json:"recommended_zinc_mg_week,omitempty"`
	RecommendedCopperMgWeek          *float64   `json:"recommended_copper_mg_week,omitempty"`
	RecommendedManganeseMgWeek       *float64   `json:"recommended_manganese_mg_week,omitempty"`
	RecommendedSeleniumMcgWeek       *float64   `json:"recommended_selenium_mcg_week,omitempty"`
	RecommendedVitaminAMcgWeek       *float64   `json:"recommended_vitamin_a_mcg_week,omitempty"`
	RecommendedVitaminCMgWeek        *float64   `json:"recommended_vitamin_c_mg_week,omitempty"`
	RecommendedVitaminDMcgWeek       *float64   `json:"recommended_vitamin_d_mcg_week,omitempty"`
	RecommendedVitaminEMgWeek        *float64   `json:"recommended_vitamin_e_mg_week,omitempty"`
	RecommendedVitaminKMcgWeek       *float64   `json:"recommended_vitamin_k_mcg_week,omitempty"`
	RecommendedThiaminMgWeek         *float64   `json:"recommended_thiamin_mg_week,omitempty"`
	RecommendedRiboflavinMgWeek      *float64   `json:"recommended_riboflavin_mg_week,omitempty"`
	RecommendedNiacinMgWeek          *float64   `json:"recommended_niacin_mg_week,omitempty"`
	RecommendedPantothenicAcidMgWeek *float64   `json:"recommended_pantothenic_acid_mg_week,omitempty"`
	RecommendedVitaminB6MgWeek       *float64   `json:"recommended_vitamin_b6_mg_week,omitempty"`
	RecommendedBiotinMcgWeek         *float64   `json:"recommended_biotin_mcg_week,omitempty"`
	RecommendedFolateMcgWeek         *float64   `json:"recommended_folate_mcg_week,omitempty"`
	RecommendedVitaminB12McgWeek     *float64   `json:"recommended_vitamin_b12_mcg_week,omitempty"`
	RecommendedSodiumMgDay           *float64   `json:"recommended_sodium_mg_day,omitempty"`
	RecommendedCholesterolMgDay      *float64   `json:"recommended_cholesterol_mg_day,omitempty"`
	RecommendedCaffeineMgDay         *float64   `json:"recommended_caffeine_mg_day,omitempty"`
	CreatedAt                        time.Time  `json:"created_at"`
	UpdatedAt                        time.Time  `json:"updated_at"`
}

type ProfileUpdate struct {
	DisplayName                  *string  `json:"display_name,omitempty"`
	DateOfBirth                  *string  `json:"date_of_birth,omitempty"`
	Age                          *int32   `json:"age,omitempty"`
	HeightCm                     *float64 `json:"height_cm,omitempty"`
	Sex                          *string  `json:"sex,omitempty"`
	ActivityLevel                *string  `json:"activity_level,omitempty"`
	Timezone                     string   `json:"timezone,omitempty"`
	MaintenanceCalories          *float64 `json:"maintenance_calories,omitempty"`
	RecommendedCalories          *float64 `json:"recommended_calories,omitempty"`
	RecommendedProteinGDay       *float64 `json:"recommended_protein_g_day,omitempty"`
	RecommendedCarbohydratesGDay *float64 `json:"recommended_carbohydrates_g_day,omitempty"`
	RecommendedFatGDay           *float64 `json:"recommended_fat_g_day,omitempty"`
	RecommendedFiberGDay         *float64 `json:"recommended_fiber_g_day,omitempty"`
	RecommendedSugarsGDay        *float64 `json:"recommended_sugars_g_day,omitempty"`
	RecommendedSaturatedFatGDay  *float64 `json:"recommended_saturated_fat_g_day,omitempty"`
}

type Service struct {
	queries *sqlc.Queries
	tokens  TokenService
}

func NewService(pool *pgxpool.Pool, tokens TokenService) *Service {
	return &Service{queries: sqlc.New(pool), tokens: tokens}
}
func (s *Service) GetProfile(ctx context.Context, userID uuid.UUID) (Profile, error) {
	user, err := s.queries.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return Profile{}, err
	}
	return profileFromRow(user), nil
}
func (s *Service) GetRecommendations(ctx context.Context, userID uuid.UUID) (Profile, error) {
	return s.GetProfile(ctx, userID)
}
func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, input ProfileUpdate) (Profile, error) {
	row, err := s.queries.GetUserByID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return Profile{}, err
	}
	profile := profileFromRow(row)
	if input.DisplayName != nil {
		profile.DisplayName = input.DisplayName
	}
	if input.DateOfBirth != nil {
		parsed, err := time.Parse(time.RFC3339, *input.DateOfBirth)
		if err == nil {
			profile.DateOfBirth = &parsed
		}
	}
	if input.Age != nil {
		profile.Age = input.Age
	}
	if input.HeightCm != nil {
		profile.HeightCm = input.HeightCm
	}
	if input.Sex != nil {
		profile.Sex = input.Sex
	}
	if input.ActivityLevel != nil {
		profile.ActivityLevel = input.ActivityLevel
	}
	if input.Timezone != "" {
		profile.Timezone = input.Timezone
	}
	if profile.Timezone == "" {
		profile.Timezone = "Asia/Kolkata"
	}
	if input.MaintenanceCalories != nil {
		profile.MaintenanceCalories = input.MaintenanceCalories
	}
	if input.RecommendedCalories != nil {
		profile.RecommendedCalories = input.RecommendedCalories
	}
	if input.RecommendedProteinGDay != nil {
		profile.RecommendedProteinGDay = input.RecommendedProteinGDay
	}
	if input.RecommendedCarbohydratesGDay != nil {
		profile.RecommendedCarbohydratesGDay = input.RecommendedCarbohydratesGDay
	}
	if input.RecommendedFatGDay != nil {
		profile.RecommendedFatGDay = input.RecommendedFatGDay
	}
	if input.RecommendedFiberGDay != nil {
		profile.RecommendedFiberGDay = input.RecommendedFiberGDay
	}
	if input.RecommendedSugarsGDay != nil {
		profile.RecommendedSugarsGDay = input.RecommendedSugarsGDay
	}
	if input.RecommendedSaturatedFatGDay != nil {
		profile.RecommendedSaturatedFatGDay = input.RecommendedSaturatedFatGDay
	}
	updated, err := s.queries.UpdateUserProfile(ctx,
		sqlc.UpdateUserProfileParams{ID: pgtype.UUID{Bytes: userID, Valid: true}, DisplayName: profile.DisplayName, DateOfBirth: dateValue(profile.DateOfBirth), Age: int32Value(profile.Age), HeightCm: numericValue(profile.HeightCm), Sex: profile.Sex, ActivityLevel: profile.ActivityLevel, Column8: profile.Timezone, MaintenanceCalories: numericValue(profile.MaintenanceCalories), RecommendedCalories: numericValue(profile.RecommendedCalories), RecommendedProteinGDay: numericValue(profile.RecommendedProteinGDay), RecommendedCarbohydratesGDay: numericValue(profile.RecommendedCarbohydratesGDay), RecommendedFatGDay: numericValue(profile.RecommendedFatGDay), RecommendedFiberGDay: numericValue(profile.RecommendedFiberGDay), RecommendedSugarsGDay: numericValue(profile.RecommendedSugarsGDay), RecommendedSaturatedFatGDay: numericValue(profile.RecommendedSaturatedFatGDay), RecommendedCalciumMgWeek: numericValue(profile.RecommendedCalciumMgWeek), RecommendedIronMgWeek: numericValue(profile.RecommendedIronMgWeek), RecommendedMagnesiumMgWeek: numericValue(profile.RecommendedMagnesiumMgWeek), RecommendedPhosphorusMgWeek: numericValue(profile.RecommendedPhosphorusMgWeek), RecommendedPotassiumMgWeek: numericValue(profile.RecommendedPotassiumMgWeek), RecommendedZincMgWeek: numericValue(profile.RecommendedZincMgWeek), RecommendedCopperMgWeek: numericValue(profile.RecommendedCopperMgWeek), RecommendedManganeseMgWeek: numericValue(profile.RecommendedManganeseMgWeek), RecommendedSeleniumMcgWeek: numericValue(profile.RecommendedSeleniumMcgWeek), RecommendedVitaminAMcgWeek: numericValue(profile.RecommendedVitaminAMcgWeek), RecommendedVitaminCMgWeek: numericValue(profile.RecommendedVitaminCMgWeek), RecommendedVitaminDMcgWeek: numericValue(profile.RecommendedVitaminDMcgWeek), RecommendedVitaminEMgWeek: numericValue(profile.RecommendedVitaminEMgWeek), RecommendedVitaminKMcgWeek: numericValue(profile.RecommendedVitaminKMcgWeek), RecommendedThiaminMgWeek: numericValue(profile.RecommendedThiaminMgWeek), RecommendedRiboflavinMgWeek: numericValue(profile.RecommendedRiboflavinMgWeek), RecommendedNiacinMgWeek: numericValue(profile.RecommendedNiacinMgWeek), RecommendedPantothenicAcidMgWeek: numericValue(profile.RecommendedPantothenicAcidMgWeek), RecommendedVitaminB6MgWeek: numericValue(profile.RecommendedVitaminB6MgWeek), RecommendedBiotinMcgWeek: numericValue(profile.RecommendedBiotinMcgWeek), RecommendedFolateMcgWeek: numericValue(profile.RecommendedFolateMcgWeek), RecommendedVitaminB12McgWeek: numericValue(profile.RecommendedVitaminB12McgWeek), RecommendedSodiumMgDay: numericValue(profile.RecommendedSodiumMgDay), RecommendedCholesterolMgDay: numericValue(profile.RecommendedCholesterolMgDay), RecommendedCaffeineMgDay: numericValue(profile.RecommendedCaffeineMgDay)})
	if err != nil {
		return Profile{}, err
	}
	return profileFromRow(updated), nil
}
func profileFromRow(user sqlc.User) Profile {
	return Profile{
		ID:                               uuid.UUID(user.ID.Bytes),
		Email:                            user.Email,
		DisplayName:                      stringPtr(user.DisplayName),
		DateOfBirth:                      timePtr(user.DateOfBirth),
		Age:                              int32Ptr(user.Age),
		HeightCm:                         float64Ptr(user.HeightCm),
		Sex:                              stringPtr(user.Sex),
		ActivityLevel:                    stringPtr(user.ActivityLevel),
		Timezone:                         user.Timezone,
		MaintenanceCalories:              float64Ptr(user.MaintenanceCalories),
		RecommendedCalories:              float64Ptr(user.RecommendedCalories),
		RecommendedProteinGDay:           float64Ptr(user.RecommendedProteinGDay),
		RecommendedCarbohydratesGDay:     float64Ptr(user.RecommendedCarbohydratesGDay),
		RecommendedFatGDay:               float64Ptr(user.RecommendedFatGDay),
		RecommendedFiberGDay:             float64Ptr(user.RecommendedFiberGDay),
		RecommendedSugarsGDay:            float64Ptr(user.RecommendedSugarsGDay),
		RecommendedSaturatedFatGDay:      float64Ptr(user.RecommendedSaturatedFatGDay),
		RecommendedCalciumMgWeek:         float64Ptr(user.RecommendedCalciumMgWeek),
		RecommendedIronMgWeek:            float64Ptr(user.RecommendedIronMgWeek),
		RecommendedMagnesiumMgWeek:       float64Ptr(user.RecommendedMagnesiumMgWeek),
		RecommendedPhosphorusMgWeek:      float64Ptr(user.RecommendedPhosphorusMgWeek),
		RecommendedPotassiumMgWeek:       float64Ptr(user.RecommendedPotassiumMgWeek),
		RecommendedZincMgWeek:            float64Ptr(user.RecommendedZincMgWeek),
		RecommendedCopperMgWeek:          float64Ptr(user.RecommendedCopperMgWeek),
		RecommendedManganeseMgWeek:       float64Ptr(user.RecommendedManganeseMgWeek),
		RecommendedSeleniumMcgWeek:       float64Ptr(user.RecommendedSeleniumMcgWeek),
		RecommendedVitaminAMcgWeek:       float64Ptr(user.RecommendedVitaminAMcgWeek),
		RecommendedVitaminCMgWeek:        float64Ptr(user.RecommendedVitaminCMgWeek),
		RecommendedVitaminDMcgWeek:       float64Ptr(user.RecommendedVitaminDMcgWeek),
		RecommendedVitaminEMgWeek:        float64Ptr(user.RecommendedVitaminEMgWeek),
		RecommendedVitaminKMcgWeek:       float64Ptr(user.RecommendedVitaminKMcgWeek),
		RecommendedThiaminMgWeek:         float64Ptr(user.RecommendedThiaminMgWeek),
		RecommendedRiboflavinMgWeek:      float64Ptr(user.RecommendedRiboflavinMgWeek),
		RecommendedNiacinMgWeek:          float64Ptr(user.RecommendedNiacinMgWeek),
		RecommendedPantothenicAcidMgWeek: float64Ptr(user.RecommendedPantothenicAcidMgWeek),
		RecommendedVitaminB6MgWeek:       float64Ptr(user.RecommendedVitaminB6MgWeek),
		RecommendedBiotinMcgWeek:         float64Ptr(user.RecommendedBiotinMcgWeek),
		RecommendedFolateMcgWeek:         float64Ptr(user.RecommendedFolateMcgWeek),
		RecommendedVitaminB12McgWeek:     float64Ptr(user.RecommendedVitaminB12McgWeek),
		RecommendedSodiumMgDay:           float64Ptr(user.RecommendedSodiumMgDay),
		RecommendedCholesterolMgDay:      float64Ptr(user.RecommendedCholesterolMgDay),
		RecommendedCaffeineMgDay:         float64Ptr(user.RecommendedCaffeineMgDay),
		CreatedAt:                        user.CreatedAt.Time,
		UpdatedAt:                        user.UpdatedAt.Time,
	}
}
func numericValue(value *float64) pgtype.Numeric {
	if value == nil {
		return pgtype.Numeric{Valid: false}
	}
	return pgtype.Numeric{Int: newIntFromFloat(*value), Exp: -3, Valid: true}
}
func int32Value(value *int32) *int32 {
	return value
}
func dateValue(value *time.Time) pgtype.Date {
	if value == nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: *value, Valid: true}
}
func float64Ptr(value pgtype.Numeric) *float64 {
	if !value.Valid {
		return nil
	}
	result, err := value.Float64Value()
	if err != nil || !result.Valid {
		return nil
	}
	v := result.Float64
	return &v
}
func int32Ptr(value *int32) *int32 {
	return value
}
func timePtr(value pgtype.Date) *time.Time {
	if !value.Valid {
		return nil
	}
	v := value.Time
	return &v
}
func stringPtr(value *string) *string        { return value }
func newIntFromFloat(value float64) *big.Int { return big.NewInt(int64(value * 1000)) }
func (s *Service) Register(ctx context.Context, email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if !strings.Contains(email, "@") || len(password) < 8 {
		return "", ErrInvalidCredentials
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{ID: pgtype.UUID{Bytes: uuid.New(), Valid: true}, Email: email, PasswordHash: string(hash)})
	if err != nil {
		return "", err
	}
	return s.tokens.Issue(uuid.UUID(user.ID.Bytes))
}
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.queries.GetUserByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", ErrInvalidCredentials
	}
	return s.tokens.Issue(uuid.UUID(user.ID.Bytes))
}
