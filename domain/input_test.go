package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"accounts/util"
)

type testcaseAccount struct {
	Input AccountInput

	Person    PersonTable
	Account   AccountTable
	Likes     []LikeTable
	Interests []InterestTable
	City      *CityTable
	Country   *CountryTable
}

var (
	now                  = time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local).Unix()
	goodAccountTestcases = []testcaseAccount{
		{
			Input: AccountInput{
				ID:     (*FieldID)(util.PtrInt32(1)),
				Email:  (*FieldEmail)(util.PtrString("test1@test.ru")),
				Sex:    (*FieldSex)(util.PtrString("m")),
				Birth:  (*FieldBirth)(util.PtrInt64(maxBirth)),
				Joined: (*FieldJoined)(util.PtrInt64(1485724260)),
				Status: (*FieldStatus)(util.PtrString("свободны")),
			},
			Person: PersonTable{
				ID:    1,
				Email: "test1@test.ru",
				Sex:   "m",
				Birth: maxBirth,
			},
			Account: AccountTable{
				ID:     1,
				Joined: 1485724260,
				Status: "свободны",
			},
		},
		{
			Input: AccountInput{
				ID:      (*FieldID)(util.PtrInt32(1)),
				Email:   (*FieldEmail)(util.PtrString("test1@test.ru")),
				Sex:     (*FieldSex)(util.PtrString("m")),
				Birth:   (*FieldBirth)(util.PtrInt64(maxBirth)),
				Joined:  (*FieldJoined)(util.PtrInt64(1485724260)),
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
					Start: (*FieldPremium)(util.PtrInt64(now)),
					End:   (*FieldPremium)(util.PtrInt64(now)),
				},
				Likes: []*AccountLikeInput{
					{
						UserID:    (*FieldID)(util.PtrInt32(1)),
						Timestamp: (*FieldTimestamp)(util.PtrInt64(now)),
					},
					{
						UserID:    (*FieldID)(util.PtrInt32(2)),
						Timestamp: (*FieldTimestamp)(util.PtrInt64(now)),
					},
				},
			},
			Person: PersonTable{
				ID:      1,
				Email:   "test1@test.ru",
				Sex:     "m",
				Birth:   maxBirth,
				Name:    util.PtrString("Андрей"),
				Surname: util.PtrString("Горбик"),
				Phone:   util.PtrString("8(999)7654321"),
			},
			Account: AccountTable{
				ID:           1,
				Joined:       1485724260,
				Status:       "заняты",
				PremiumStart: util.PtrInt64(now),
				PremiumEnd:   util.PtrInt64(now),
			},
			Likes: []LikeTable{
				{
					LikerID:   1,
					LikeeID:   1,
					Timestamp: now,
				},
				{
					LikerID:   1,
					LikeeID:   2,
					Timestamp: now,
				},
			},
			Interests: []InterestTable{
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
			City: &CityTable{
				Name: "Москва",
			},
			Country: &CountryTable{
				Name: "Россия",
			},
		},
	}

	badAccountTestcases = []testcaseAccount{}
)

func Test_AccountInputToTables_Success(t *testing.T) {
	for _, testcase := range goodAccountTestcases {
		if err := testcase.Input.Validate(); err != nil {
			t.Error(err)
			continue
		}

		assert.Equal(t, testcase.Account, testcase.Input.GetAccount())
		assert.Equal(t, testcase.Person, testcase.Input.GetPerson())
		assert.Equal(t, testcase.Interests, testcase.Input.GetInterests())
		assert.Equal(t, testcase.Likes, testcase.Input.GetLikes())
		assert.Equal(t, testcase.City, testcase.Input.GetCity())
		assert.Equal(t, testcase.Country, testcase.Input.GetCountry())
	}
}
