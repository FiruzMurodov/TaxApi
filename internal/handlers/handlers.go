package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"taxApi/internal/logs"
	"taxApi/internal/models"
	"taxApi/internal/services"
)

type Handler struct {
	Service *services.Service
	Lg      *logs.Logger
}

func NewHandler(service *services.Service, lg *logs.Logger) *Handler {
	return &Handler{Service: service, Lg: lg}
}

func (h *Handler) Report(w http.ResponseWriter, r *http.Request) {
	var reqReport *models.ReqReport
	err := json.NewDecoder(r.Body).Decode(&reqReport)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ConvertToJson(models.ErrInvalidData))
		h.Lg.Error()
		return
	}

	err = h.Service.Report(reqReport)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ConvertToJson(models.ErrInternal))
		h.Lg.Error()
		return
	}

	_, err = w.Write(ConvertToJson(models.ReportStatus))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *Handler) ReportDriver(w http.ResponseWriter, r *http.Request) {
	var reqReport *models.ReqReport
	err := json.NewDecoder(r.Body).Decode(&reqReport)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ConvertToJson(models.ErrInvalidData))
		h.Lg.Error()
		return
	}
	idUser, _ := h.Service.IdByToken(r.Context(), r.Header.Get("Authorization"))

	err, file := h.Service.ReportDriver(reqReport, idUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ConvertToJson(models.ErrInternal))
		h.Lg.Error()
		return
	}

	filePath := "report.xlsx"
	err = file.Save(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ConvertToJson(models.ErrInternal))
		h.Lg.Error()
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=report.xlsx")
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeFile(w, r, filePath)

	//w.Header().Set("Content-Type", "application/octet-stream")
	//w.Header().Set("Content-Type", "application/csv;name=report.xlsx")
	//w.Header().Set("Content-Encoding", "br")
	//w.Header().Set("Content-Transfer-Encoding", "base64")

	err = os.Remove(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ConvertToJson(models.ErrInternal))
		h.Lg.Error()
		return
	}
}

func (h *Handler) ReadById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Lg.Error()
		w.Write(ConvertToJson(models.ErrInvalidData))
		return
	}

	reportsDriver, reportsCustomer, err := h.Service.ReadOrderById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		h.Lg.Error()
		w.Write(ConvertToJson(models.ErrNotFound))
		return
	}

	if len(reportsDriver) != 0 {
		bytes, err := json.Marshal(reportsDriver)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Lg.Error()
			w.Write(ConvertToJson(models.ErrInternal))
			return
		}

		_, err = w.Write(bytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Lg.Error()
			w.Write(ConvertToJson(models.ErrInternal))
			return
		}
	}

	if len(reportsCustomer) != 0 {
		bytes, err := json.Marshal(reportsCustomer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Lg.Error()
			w.Write(ConvertToJson(models.ErrInternal))
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.Lg.Error()
			w.Write(ConvertToJson(models.ErrInternal))
			return
		}
	}

}

/*
	func (h *Handler) CreateTravel(w http.ResponseWriter, r *http.Request) {
		var reqBody *models.RequestBody
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(ConvertToJson(models.ErrInvalidData))
			h.Lg.Error()
			return
		}

		idUser, _ := h.Service.IdByToken(r.Context(), r.Header.Get("Authorization"))

		_, price, err := h.Service.GetRoute(reqBody, idUser)
		if err != nil {
			h.Lg.Error()
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(ConvertToJson(models.ErrInternal))
			return
		}
		respBody := models.OkStatus + fmt.Sprintf("%.2f", price) + " сомони"

		_, err = w.Write(ConvertToJson(respBody))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//resp, err := http.Get("https://router.project-osrm.org/table/v1/driving/" + reqBody.LongitudeSource + "," + reqBody.LatitudeSource + ";" + reqBody.LongitudeDestination + "," + reqBody.LatitudeDestination + "?sources=1&destinations=0&annotations=duration,distance")
		//resp, err := http.Get(`https://router.project-osrm.org/table/v1/car/55.745163,37.624513;55.738450,37.619319?sources=1&destinations=0&annotations=duration,distance`)
		//resp, err := http.Get("http://router.project-osrm.org/tile/v1/car/tile(1310,3166,13).mvt")
		//resp, err := http.Get("http://router.project-osrm.org/table/v1/car/" + reqBody.LongitudeSource + "," + reqBody.LatitudeSource + ";" + reqBody.LongitudeDestination + "," + reqBody.LatitudeDestination + "?source=0;1")
		//resp, err := http.Get("http://router.project-osrm.org/trip/v1/driving/55.745163,37.624513;55.738450,37.619319?source=first&destination=last&annotations=duration,distance")
		//resp, err := http.Get("http://localhost:9999/trip/v1/driving/48.565164,13.431038;48.574931,13.465600")

		//if err != nil {
		//	h.Lg.Warn("oshibka seti")
		//	log.Fatal(err)
		//}

		//log.Println(resp.Location())
		//var byte1 *http.Response
		//_, err = json.Marshal(body)
		//log.Println(body)
		//fmt.Printf("%s", body)
	}
*/
func (h *Handler) CreateTravel(w http.ResponseWriter, r *http.Request) {
	var reqBody *models.RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ConvertToJson(models.ErrInvalidData))
		h.Lg.Error()
		return
	}

	idUser, _ := h.Service.IdByToken(r.Context(), r.Header.Get("Authorization"))

	route, price, err := h.Service.GetRoute(reqBody, idUser)
	if err != nil {
		h.Lg.Error()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ConvertToJson(models.ErrInternal))
		return
	}

	message := models.OkStatus + fmt.Sprintf("%.2f", price) + " сомони"
	routeDuration := int(route.Durations[0][0])
	routeDurationMinute := routeDuration / 60
	routeDurationSecond := routeDuration % 60
	duration := strconv.Itoa(routeDurationMinute) + "мин. " + strconv.Itoa(routeDurationSecond) + " c."

	bytes, err := json.Marshal(models.MessageToCustomer{
		Message:  message,
		Duration: duration,
	})
	if err != nil {
		h.Lg.Error()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ConvertToJson(models.ErrInternal))
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		h.Lg.Error()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ConvertToJson(models.ErrInternal))
		return
	}
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idUser, _ := h.Service.IdByToken(r.Context(), r.Header.Get("Authorization"))

	orderToCustomer, err := h.Service.GetTravelToCustomer(idUser)
	if err != nil {
		h.Lg.Error()
		w.WriteHeader(http.StatusNotFound)
		w.Write(ConvertToJson(models.ErrNotOrders))
		return
	}

	bytes, err := json.Marshal(orderToCustomer)
	if err != nil {
		h.Lg.Error()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ConvertToJson(models.ErrInternal))
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		h.Lg.Error()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ConvertToJson(models.ErrInternal))
		return
	}
}

func (h *Handler) GetTravelList(w http.ResponseWriter, r *http.Request) {

	idUser, _ := h.Service.IdByToken(r.Context(), r.Header.Get("Authorization"))

	orderToDriver, err := h.Service.GetTravelToDriver(idUser)
	if err != nil {
		h.Lg.Error()
		w.WriteHeader(http.StatusNotFound)
		w.Write(ConvertToJson(models.ErrNotOrders))
		return
	}
	bytes, err := json.Marshal(orderToDriver)
	if err != nil {
		w.Write(ConvertToJson(models.ErrInternal))
		w.WriteHeader(http.StatusInternalServerError)
		h.Lg.Error()
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		w.Write(ConvertToJson(models.ErrInternal))
		w.WriteHeader(http.StatusInternalServerError)
		h.Lg.Error()
		return
	}

}

func (h *Handler) GetTravelById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Lg.Error()
		w.Write(ConvertToJson(models.ErrInvalidData))
		return
	}

	idUser, _ := h.Service.IdByToken(r.Context(), r.Header.Get("Authorization"))

	orderToDriver, err := h.Service.GetTravelToDriver(idUser)
	if err != nil {
		h.Lg.Error()
		w.WriteHeader(http.StatusNotFound)
		w.Write(ConvertToJson(models.ErrNotOrders))
		return
	}

	matchId := 0
	for _, driver := range orderToDriver {
		if driver.Id == id {
			matchId++
			break
		}
	}
	if matchId != 1 {
		h.Lg.Error()
		w.WriteHeader(http.StatusNotFound)
		w.Write(ConvertToJson(models.ErrNotValidateOrder))
		return
	}

	travelById, err := h.Service.GetTravelById(orderToDriver, id, idUser)
	if err != nil {
		h.Lg.Error()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ConvertToJson(models.ErrInternal))
		return
	}
	orderResponse := models.OrderResponseToDriver{
		Message:       models.OrderSuccess + strconv.Itoa(travelById.Id),
		CustomerName:  travelById.CustomerName,
		CustomerPhone: travelById.CustomerPhone,
		Price:         travelById.Price,
	}

	bytes, err := json.Marshal(orderResponse)
	if err != nil {
		w.Write(ConvertToJson(models.ErrInternal))
		w.WriteHeader(http.StatusInternalServerError)
		h.Lg.Error()
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		w.Write(ConvertToJson(models.ErrInternal))
		w.WriteHeader(http.StatusInternalServerError)
		h.Lg.Error()
		return
	}

}

func (h *Handler) EndTravelById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idTravel, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Lg.Error()
		w.Write(ConvertToJson(models.ErrInvalidData))
		return
	}

	idUser, _ := h.Service.IdByToken(r.Context(), r.Header.Get("Authorization"))

	err = h.Service.EndTravelById(idUser, idTravel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.Lg.Error()
		w.Write(ConvertToJson(models.ErrNotFoundOrder))
		return
	}

	_, err = w.Write(ConvertToJson(models.SuccessOrder))
	if err != nil {
		w.Write(ConvertToJson(models.ErrInternal))
		w.WriteHeader(http.StatusInternalServerError)
		h.Lg.Error()
		return
	}

}

func (h *Handler) GetTokenToUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ConvertToJson(models.ErrInvalidData))
		h.Lg.Error()
		return
	}

	tokenToUser, err := h.Service.GetTokenToUser(r.Context(), user)
	if err != nil {
		h.Lg.Error()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ConvertToJson(models.ErrInvalidLoginOrPassword))
		return
	}

	item := models.Token{TokenText: tokenToUser}

	data, err := json.Marshal(item)
	if err != nil {
		w.Write(ConvertToJson(models.ErrInternal))
		w.WriteHeader(http.StatusInternalServerError)
		h.Lg.Error()
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		w.Write(ConvertToJson(models.ErrInternal))
		w.WriteHeader(http.StatusInternalServerError)
		h.Lg.Error()
		return
	}
}
