package pkg

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	ErrBadBodyRequest                      = errors.New("bad body request")
	ErrJSONUnexpectedEnd                   = errors.New("unexpected end of JSON input")
	ErrContentTypeUndefined                = errors.New("content-type undefined")
	ErrUnsupportedMediaType                = errors.New("unsupported media type")
	ErrEmptyBody                           = errors.New("empty body")
	ErrConvertQueryType                    = errors.New("bad input query")
	ErrQueryRequiredEmpty                  = errors.New("miss query params")
	ErrBadRequestParams                    = errors.New("bad query params")
	ErrBadRequestParamsEmptyRequiredFields = errors.New("bad params, empty required field")
	ErrGetEasyJSON                         = errors.New("err get easyjson")

	ErrNotFoundInDB             = errors.New("not found")
	ErrWorkDatabase             = errors.New("error sql")
	ErrGetParamsConvert         = errors.New("err get sql params")
	ErrUnsupportedSortParameter = errors.New("unsupported sort parameter")

	ErrBigRequest    = errors.New("big request")
	ErrConvertLength = errors.New("getting content-length failed")

	ErrSuchUserExist          = errors.New("such user exist")
	ErrSuchUserNotFound       = errors.New("such user not fount")
	ErrUpdateUserDataConflict = errors.New("impossible update such user data")

	ErrSuchThreadNotFound = errors.New("such thread not fount")
	ErrSuchThreadExist    = errors.New("such thread exist")

	ErrNoSuchRuleSortPosts = errors.New("no such rule for sort posts")
	ErrSuchPostNotFound    = errors.New("such post not found")
	ErrPostParentNotFound  = errors.New("such post parent not found")
	ErrInvalidParent       = errors.New("parent not valid")

	ErrSuchForumNotFound = errors.New("such forum not fount")
	ErrSuchForumExist    = errors.New("such forum exist")
)

type ErrHTTPClassifier struct {
	table map[string]int
}

func NewErrHTTPClassifier() ErrHTTPClassifier {
	res := make(map[string]int)

	res[ErrBadBodyRequest.Error()] = http.StatusBadRequest
	res[ErrJSONUnexpectedEnd.Error()] = http.StatusBadRequest
	res[ErrContentTypeUndefined.Error()] = http.StatusBadRequest
	res[ErrUnsupportedMediaType.Error()] = http.StatusUnsupportedMediaType
	res[ErrEmptyBody.Error()] = http.StatusBadRequest
	res[ErrConvertQueryType.Error()] = http.StatusBadRequest
	res[ErrQueryRequiredEmpty.Error()] = http.StatusBadRequest
	res[ErrBadRequestParams.Error()] = http.StatusBadRequest
	res[ErrBadRequestParamsEmptyRequiredFields.Error()] = http.StatusBadRequest
	res[ErrBadRequestParams.Error()] = http.StatusBadRequest
	res[ErrGetEasyJSON.Error()] = http.StatusInternalServerError

	res[ErrNotFoundInDB.Error()] = http.StatusNotFound
	res[ErrWorkDatabase.Error()] = http.StatusInternalServerError
	res[ErrGetParamsConvert.Error()] = http.StatusInternalServerError
	res[ErrUnsupportedSortParameter.Error()] = http.StatusBadRequest

	res[ErrBigRequest.Error()] = http.StatusBadRequest
	res[ErrConvertLength.Error()] = http.StatusBadRequest

	res[ErrSuchUserExist.Error()] = http.StatusConflict
	res[ErrSuchUserNotFound.Error()] = http.StatusNotFound
	res[ErrUpdateUserDataConflict.Error()] = http.StatusConflict

	res[ErrSuchThreadNotFound.Error()] = http.StatusNotFound

	res[ErrNoSuchRuleSortPosts.Error()] = http.StatusNotFound
	res[ErrSuchPostNotFound.Error()] = http.StatusNotFound
	res[ErrPostParentNotFound.Error()] = http.StatusNotFound
	res[ErrPostParentNotFound.Error()] = http.StatusConflict

	res[ErrSuchForumNotFound.Error()] = http.StatusNotFound
	res[ErrInvalidParent.Error()] = http.StatusConflict

	return ErrHTTPClassifier{
		table: res,
	}
}

var errHTTPCsf = NewErrHTTPClassifier()

func GetErrorCodeHTTP(err error) (int, bool) {
	code, exist := errHTTPCsf.table[err.Error()]
	if !exist {
		return http.StatusInternalServerError, exist
	}

	return code, exist
}
