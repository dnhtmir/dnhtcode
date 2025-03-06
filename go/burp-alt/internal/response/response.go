package response

import (
	"fmt"
	"io"
	"net/http"
)

func HandleResponse(response *http.Response) {
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response Status:", response.Status)
	fmt.Println("Response Headers:", response.Header)
	fmt.Println("Response Body:", string(body))
}
