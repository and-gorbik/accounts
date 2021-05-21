package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"accounts/util"
)

type testcaseAccount struct {
	Input AccountInput

	PersonModel    PersonModel
	AccountModel   AccountModel
	LikeModels     []LikeModel
	InterestModels []InterestModel
	CityModel      *CityModel
	CountryModel   *CountryModel
}

var (
	testNow              = time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local)
	testBirth            = time.Date(1994, 1, 1, 0, 0, 0, 0, time.Local)
	goodAccountTestcases = []testcaseAccount{
		{
			Input: AccountInput{
				ID:     (*FieldID)(util.PtrInt32(1)),
				Email:  (*FieldEmail)(util.PtrString("test1@test.ru")),
				Sex:    (*FieldSex)(util.PtrString("m")),
				Birth:  (*FieldBirth)(util.PtrInt64(testBirth.Unix())),
				Joined: (*FieldJoined)(util.PtrInt64(testNow.Unix())),
				Status: (*FieldStatus)(util.PtrString("свободны")),
			},
			PersonModel: PersonModel{
				ID:     1,
				Status: "свободны",
				Email:  "test1@test.ru",
				Sex:    "m",
				Birth:  testBirth,
			},
			AccountModel: AccountModel{
				ID:     1,
				Joined: testNow,
			},
		},
		{
			Input: AccountInput{
				ID:      (*FieldID)(util.PtrInt32(1)),
				Email:   (*FieldEmail)(util.PtrString("test1@test.ru")),
				Sex:     (*FieldSex)(util.PtrString("m")),
				Birth:   (*FieldBirth)(util.PtrInt64(testBirth.Unix())),
				Joined:  (*FieldJoined)(util.PtrInt64(testNow.Unix())),
				Status:  (*FieldStatus)(util.PtrString("заняты")),
				Name:    (*FieldFirstname)(util.PtrString("Андрей")),
				Surname: (*FieldSurname)(util.PtrString("Горбик")),
				Phone:   (*FieldPhone)(util.PtrString("8(999)7654321")),
				Country: (*FieldCountry)(util.PtrString("Россия")),
				City:    (*FieldCity)(util.PtrString("Москва")),
				Interests: []*FieldInterest{
					(*FieldInterest)(util.PtrString("компьютер")),
					(*FieldInterest)(util.PtrString("волейбол")),
					(*FieldInterest)(util.PtrString("фортепиано")),
				},
				Premium: &PremiumInput{
					Start: (*FieldPremium)(util.PtrInt64(testNow.Unix())),
					End:   (*FieldPremium)(util.PtrInt64(testNow.Unix())),
				},
				Likes: []*AccountLikeInput{
					{
						UserID:    (*FieldID)(util.PtrInt32(1)),
						Timestamp: (*FieldTimestamp)(util.PtrInt64(testNow.Unix())),
					},
					{
						UserID:    (*FieldID)(util.PtrInt32(2)),
						Timestamp: (*FieldTimestamp)(util.PtrInt64(testNow.Unix())),
					},
				},
			},
			PersonModel: PersonModel{
				ID:      1,
				Status:  "заняты",
				Email:   "test1@test.ru",
				Sex:     "m",
				Birth:   testBirth,
				Name:    util.PtrString("Андрей"),
				Surname: util.PtrString("Горбик"),
				Phone:   util.PtrString("8(999)7654321"),
			},
			AccountModel: AccountModel{
				ID:           1,
				Joined:       testNow,
				PremiumStart: &testNow,
				PremiumEnd:   &testNow,
			},
			LikeModels: []LikeModel{
				{
					LikerID:   1,
					LikeeID:   1,
					Timestamp: testNow,
				},
				{
					LikerID:   1,
					LikeeID:   2,
					Timestamp: testNow,
				},
			},
			InterestModels: []InterestModel{
				{
					AccountID: 1,
					Name:      "компьютер",
				},
				{
					AccountID: 1,
					Name:      "волейбол",
				},
				{
					AccountID: 1,
					Name:      "фортепиано",
				},
			},
			CityModel: &CityModel{
				Name: "Москва",
			},
			CountryModel: &CountryModel{
				Name: "Россия",
			},
		},
	}

	badAccountTestcases = []testcaseAccount{}
)

func Test_AccountInputToModels_Success(t *testing.T) {
	for _, testcase := range goodAccountTestcases {
		if err := testcase.Input.Validate(); err != nil {
			t.Error(err)
			continue
		}

		assert.Equal(t, testcase.AccountModel, *testcase.Input.AccountModel())
		assert.Equal(t, testcase.PersonModel, *testcase.Input.PersonModel(nil, nil))
		assert.Equal(t, testcase.InterestModels, testcase.Input.InterestModels())
		assert.Equal(t, testcase.LikeModels, testcase.Input.LikeModels())
		assert.Equal(t, testcase.CityModel, testcase.Input.CityModel())
		assert.Equal(t, testcase.CountryModel, testcase.Input.CountryModel())
	}
}
