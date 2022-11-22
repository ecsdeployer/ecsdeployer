package helpers_test

import (
	"fmt"
	"math/rand"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/require"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func TestNameValuePairMerger(t *testing.T) {

	pairMapSize := 10

	pairMap := make(map[string]string, pairMapSize)
	for i := 1; i <= pairMapSize; i++ {
		pairMap[fmt.Sprintf("NAME%02d", i)] = randSeq(10)
	}

	pairMap2 := make(map[string]string, pairMapSize)
	for i := 1; i <= pairMapSize; i++ {
		pairMap2[fmt.Sprintf("NAME%02d", i)] = randSeq(10)
	}

	pairGroup1 := buildPairs(pairMap, "NAME01", "NAME02")
	pairGroup2 := buildPairs(pairMap, "NAME03", "NAME04")

	pairGroup3 := buildPairs(pairMap2, "NAME01")
	pairGroup4 := buildPairs(pairMap2, "NAME03")

	pairGroup5 := buildPairs(pairMap, "NAME05", "NAME06")

	t.Run("no dupes, no blanks", func(t *testing.T) {
		merged := helpers.NameValuePairMerger(pairGroup1, nil, pairGroup2, pairGroup5)
		require.Truef(t, ensureUniquePairs(merged), "result has duplicate pairs!")
		cleanMap := pairToMap(merged)
		for _, k := range []string{"NAME01", "NAME02", "NAME03", "NAME04", "NAME05", "NAME06"} {
			require.Contains(t, cleanMap, k)
			require.Equal(t, cleanMap[k], pairMap[k])
		}
	})

	t.Run("no dupes, with blanks", func(t *testing.T) {
		merged := helpers.NameValuePairMerger(pairGroup1, nil, pairGroup2, pairGroup5, []config.NameValuePair{
			{
				Name:  aws.String("NAME02"),
				Value: nil,
			},
		})
		require.Truef(t, ensureUniquePairs(merged), "result has duplicate pairs!")
		cleanMap := pairToMap(merged)
		for _, k := range []string{"NAME01", "NAME02", "NAME03", "NAME04", "NAME05", "NAME06"} {
			require.Contains(t, cleanMap, k)
			require.Equal(t, cleanMap[k], pairMap[k])
		}
	})

	t.Run("when duplicate pairs", func(t *testing.T) {
		merged := helpers.NameValuePairMerger(pairGroup1, pairGroup2, nil, pairGroup3, pairGroup4, pairGroup5)

		require.Truef(t, ensureUniquePairs(merged), "result has duplicate pairs!")
		cleanMap := pairToMap(merged)
		for _, k := range []string{"NAME02", "NAME04", "NAME05", "NAME06"} {
			require.Contains(t, cleanMap, k)
			require.Equal(t, cleanMap[k], pairMap[k])
		}

		for _, k := range []string{"NAME01", "NAME03"} {
			require.Contains(t, cleanMap, k)
			require.Equal(t, cleanMap[k], pairMap2[k])
		}
	})
}

func ensureUniquePairs(pairs []config.NameValuePair) bool {
	pairCounts := make(map[string]struct{}, len(pairs))

	for _, pair := range pairs {
		_, ok := pairCounts[*pair.Name]
		if ok {
			return false
		}
		pairCounts[*pair.Name] = struct{}{}
	}

	return true

}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func buildPairs(pairMap map[string]string, keys ...string) []config.NameValuePair {
	result := make([]config.NameValuePair, 0, len(keys))
	for _, k := range keys {
		v, ok := pairMap[k]
		if !ok {
			panic(fmt.Errorf("Key '%s' DOES NOT EXIST IN PAIRMAP", k))
		}
		result = append(result, config.NameValuePair{
			Name:  aws.String(k),
			Value: aws.String(v),
		})
	}
	return result
}

func pairToMap(pairs []config.NameValuePair) map[string]string {
	result := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		result[*pair.Name] = *pair.Value
	}
	return result
}
