package utils

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	hashids "github.com/speps/go-hashids/v2"
	"github.com/hoangminhphuc/goph-chat/common"
)

var (
	encoder *hashids.HashID
	once    sync.Once
	version string
)

func initEncoder() {
	data := hashids.NewData()
	seed := os.Getenv("HASHIDS_SEED")
	salt, err := GenerateFixedSalt(seed, common.DefaultSaltLength)

	if err != nil {
		panic(err)
	}

	version = os.Getenv("HASHIDS_SALT_VERSION")
	if version == "" {
		version = "v1"
	}
	data.Salt = salt
	if ml := os.Getenv("HASHIDS_MIN_LENGTH"); ml != "" {
		if n, err := strconv.Atoi(ml); err == nil {
			data.MinLength = n
		}
	}
	e, err := hashids.NewWithData(data)
	if err != nil {
		panic(common.WrapError(err, "failed to init Hashids", http.StatusInternalServerError))
	}
	encoder = e
}

func EncodeID(id uint64) (string, error) {
	once.Do(initEncoder)
	s, err := encoder.EncodeInt64([]int64{int64(id)})
	if err != nil {
		return "", common.WrapError(err, "failed to encode ID", http.StatusInternalServerError)
	}
	return version + "_" + s, nil
}

func DecodeID(hash string) (int, error) {
	once.Do(initEncoder)
	parts := strings.SplitN(hash, "_", 2)
	if len(parts) != 2 || parts[0] != version {
		return 0, common.WrapError(nil, "invalid ID format", http.StatusBadRequest)
	}
	arr, err := encoder.DecodeInt64WithError(parts[1])
	if err != nil || len(arr) == 0 {
		return 0, common.WrapError(err, "failed to decode ID", http.StatusBadRequest)
	}
	return int(arr[0]), nil
}
