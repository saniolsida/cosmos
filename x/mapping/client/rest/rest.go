package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/mapping/types"
	"github.com/gorilla/mux"
)

func RegisterRoutes(clientCtx client.Context, rtr *mux.Router) {
	rtr.HandleFunc(
		"/mapping/address/{device_id}",
		getAddressFromDeviceIDHandler(clientCtx),
	).Methods("GET")
}

func getAddressFromDeviceIDHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deviceID := mux.Vars(r)["device_id"]

		queryClient := types.NewQueryClient(clientCtx)
		res, err := queryClient.GetAddressFromDeviceID(
			r.Context(),
			&types.QueryGetAddressFromDeviceIDRequest{DeviceId: deviceID},
		)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(fmt.Sprintf("error: %s", err)))
			return
		}

		jsonBytes, err := clientCtx.Codec.MarshalInterfaceJSON(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("marshal error: %s", err)))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(jsonBytes)
	}
}
