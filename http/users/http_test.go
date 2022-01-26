package users

import (
	"user-curd/models"
	"user-curd/services"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := services.NewMockServices(ctrl)
	mock := New(mockService)

	testCases := []struct {
		desc      string
		id        string
		mock      []*gomock.Call
		expecErr  error
		expecBody []byte
	}{
		{
			desc:     "Success case ",
			id:       "2",
			expecErr: nil,
			mock: []*gomock.Call{
				mockService.EXPECT().FetchUserDetailsById(2).Return(models.User{
					Id:    2,
					Name:  "gopi",
					Email: "gopi@gmail.com",
					Phone: "1234567899",
					Age:   23,
				}, nil),
			},
			expecBody: []byte(`{"Id":2,"Name":"gopi","Email":"gopi@gmail.com","Phone":"1234567899","Age":23}`),
		},
		{
			desc:      "Failure case-1",
			id:        "1a",
			expecErr:  nil,
			expecBody: []byte("invalid parameter id"),
		},
		{
			desc:     "Failure case - 2",
			id:       "2",
			expecErr: nil,
			mock: []*gomock.Call{
				mockService.EXPECT().FetchUserDetailsById(2).Return(models.User{
					Id:    2,
					Name:  "gopi",
					Email: "gopi@gmail.com",
					Phone: "1234567899",
					Age:   23,
				}, errors.New("internal error")),
			},
			expecBody: []byte("internal error"),
		},
	}

	for _, v := range testCases {
		t.Run(v.desc, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/user?id="+v.id, nil)
			rw := httptest.NewRecorder()
			mock.GetUserById(rw, r)
			if rw.Body.String() != string(v.expecBody) {
				t.Errorf("Expected %v Obtained %v", string(v.expecBody), rw.Body.String())
			}
		})
	}
}

func TestPostUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := services.NewMockServices(ctrl)
	mock := New(mockService)

	testCases := []struct {
		desc     string
		user     models.User
		mock     []*gomock.Call
		expecErr error
		expecRes []byte
	}{
		{
			desc: "Success case",
			user: models.User{
				Id:    1,
				Name:  "gopi",
				Email: "gopi@gmail.com",
				Phone: "1234567899",
				Age:   23,
			},
			mock: []*gomock.Call{
				mockService.EXPECT().EmailValidation("gopi@gmail.com").Return(true, nil).MaxTimes(5),
				mockService.EXPECT().InsertUserDetails(models.User{
					Id:    1,
					Name:  "gopi",
					Email: "gopi@gmail.com",
					Phone: "1234567899",
					Age:   23,
				}).Return(nil).MaxTimes(5),
			},
			expecErr: nil,
			expecRes: []byte("User created"),
		},
		{
			desc: "Failure case -1",
			user: models.User{
				Id:    0,
				Name:  "gopi",
				Email: "gopi@gmail.com",
				Phone: "1234567899",
				Age:   23,
			},
			expecErr: errors.New("Id shouldn't be zero"),
			expecRes: []byte("Id shouldn't be zero"),
		},
		{
			desc: "Failure case -2",
			user: models.User{
				Id:    1,
				Name:  "gopi",
				Email: "gopi12@gmail.com",
				Phone: "1234567899",
				Age:   23,
			},
			mock: []*gomock.Call{
				mockService.EXPECT().EmailValidation("gopi12@gmail.com").Return(false, nil),
			},
			expecErr: errors.New("email already present - could not create user"),
			expecRes: []byte("email already present - could not create user"),
		},
	}
	for _, v := range testCases {
		t.Run(v.desc, func(t *testing.T) {
			b, _ := json.Marshal(v.user)
			r := httptest.NewRequest("POST", "/insert", bytes.NewBuffer(b))
			rw := httptest.NewRecorder()
			mock.PostUser(rw, r)
			fmt.Println(rw.Body.String())
			if rw.Body.String() != string(v.expecRes) {
				t.Errorf("Expected %v Obtained %v", string(v.expecRes), rw.Body.String())
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := services.NewMockServices(ctrl)
	mock := New(mockService)
	testCases := []struct {
		desc      string
		mock      []*gomock.Call
		expecErr  error
		expecBody []byte
	}{
		{
			desc:     "Success case ",
			expecErr: nil,
			mock: []*gomock.Call{
				mockService.EXPECT().FetchAllUserDetails().Return([]models.User{{
					Id:    2,
					Name:  "gopi",
					Email: "gopi@gmail.com",
					Phone: "1234567899",
					Age:   23,
				},
				}, nil),
			},
			expecBody: []byte(`[{"Id":2,"Name":"gopi","Email":"gopi@gmail.com","Phone":"1234567899","Age":23}]`),
		},
		{
			desc:     "Failure case-1 ",
			expecErr: errors.New("error generated"),
			mock: []*gomock.Call{
				mockService.EXPECT().FetchAllUserDetails().Return([]models.User{{
					Id:    2,
					Name:  "gopi",
					Email: "gopi@gmail.com",
					Phone: "1234567899",
					Age:   23,
				},
				}, errors.New("error generated")),
			},
			expecBody: []byte("error generated"),
		},
	}
	for _, v := range testCases {
		t.Run(v.desc, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/users", nil)
			rw := httptest.NewRecorder()
			mock.GetUsers(rw, r)
			if rw.Body.String() != string(v.expecBody) {
				t.Errorf("Expected %v Obtained %v", string(v.expecBody), rw.Body.String())
			}
		})
	}

}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := services.NewMockServices(ctrl)
	mock := New(mockService)
	testCases := []struct {
		desc     string
		id       string
		expecErr error
		mock     []*gomock.Call
		expecRes []byte
	}{
		{
			desc:     "Success case",
			id:       "1",
			expecErr: nil,
			mock: []*gomock.Call{
				mockService.EXPECT().DeleteUserDetailsById(1).Return(nil),
			},
			expecRes: []byte("User deleted successfully"),
		},
		{
			desc:     "Failure case - 1",
			id:       "0",
			expecErr: errors.New("Id shouldn't be zero"),
			expecRes: []byte("Id shouldn't be zero"),
		},
		{
			desc:     "Failure case - 2",
			id:       "1",
			expecErr: errors.New("error generated"),
			mock: []*gomock.Call{
				mockService.EXPECT().DeleteUserDetailsById(1).Return(errors.New("error generated")),
			},
			expecRes: []byte("error generated"),
		},
	}
	for _, v := range testCases {
		t.Run(v.desc, func(t *testing.T) {
			r := httptest.NewRequest("DELETE", "/delete?id="+v.id, nil)
			rw := httptest.NewRecorder()
			mock.DeleteUser(rw, r)
			fmt.Println(rw.Body.String())
			if rw.Body.String() != string(v.expecRes) {
				t.Errorf("Expected %v Obtained %v", string(v.expecRes), rw.Body.String())
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := services.NewMockServices(ctrl)
	mock := New(mockService)

	testCases := []struct {
		desc     string
		user     models.User
		mock     []*gomock.Call
		expecErr error
		expecRes []byte
	}{
		{
			desc: "Success case",
			user: models.User{
				Id:    1,
				Name:  "gopi",
				Email: "gopi@gmail.com",
				Phone: "1234567899",
				Age:   23,
			},
			mock: []*gomock.Call{
				mockService.EXPECT().EmailValidation("gopi@gmail.com").Return(true, nil),
				mockService.EXPECT().UpdateUserDetails(models.User{
					Id:    1,
					Name:  "gopi",
					Email: "gopi@gmail.com",
					Phone: "1234567899",
					Age:   23,
				}).Return(nil),
			},
			expecErr: nil,
			expecRes: []byte("User updated"),
		},
		{
			desc: "Failure case -1",
			user: models.User{
				Id:    0,
				Name:  "gopi",
				Email: "gopi@gmail.com",
				Phone: "1234567899",
				Age:   23,
			},
			expecErr: errors.New("Id shouldn't be zero"),
			expecRes: []byte("Id shouldn't be zero"),
		},
		{
			desc: "Failure case -2",
			user: models.User{
				Id:    1,
				Name:  "gopi",
				Email: "gopi12@gmail.com",
				Phone: "1234567899",
				Age:   23,
			},
			mock: []*gomock.Call{
				mockService.EXPECT().EmailValidation("gopi12@gmail.com").Return(false, nil),
			},
			expecErr: errors.New("email already present - could not create user"),
			expecRes: []byte("email already present - could not create user"),
		},
	}
	for _, v := range testCases {
		t.Run(v.desc, func(t *testing.T) {
			b, _ := json.Marshal(v.user)
			r := httptest.NewRequest("PUT", "/update", bytes.NewBuffer(b))
			rw := httptest.NewRecorder()
			mock.UpdateUser(rw, r)
			fmt.Println(rw.Body.String())
			if rw.Body.String() != string(v.expecRes) {
				t.Errorf("Expected %v Obtained %v", string(v.expecRes), rw.Body.String())
			}
		})
	}
}