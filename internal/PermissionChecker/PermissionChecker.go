package permissionchecker

import (
	"encoding/json"
	"fmt"

	"github.com/Blue-Onion/PNB/internal/models"
)

func PermissionChecker(res *models.Result) {
	data, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(data))
}
