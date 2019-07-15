package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/pymq/go-test-assignment/model"
	"net/http"
	"strconv"
)

func (h *Handler) GetTvs(c echo.Context) (err error) {
	var results []model.Tv
	err = h.DB.Select(&results, `SELECT * FROM "tv"`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalError())
	}
	return c.JSON(http.StatusOK, results)
}

func (h *Handler) GetTvById(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if id < 1 || err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage("Invalid id"))
	}
	tv := model.Tv{}
	err = h.DB.Get(&tv, `SELECT * 
 				FROM "tv" WHERE id=$1`, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, NotFound())
	}
	return c.JSON(http.StatusOK, tv)
}

func (h *Handler) PutTvById(c echo.Context) (err error) {
	tv := model.Tv{}
	if err = c.Bind(&tv); err != nil {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	tv.ID = int64(id)
	tv.Brand.Parse(c.FormValue("brand"))
	if !tv.IsValid() || tv.ID <= 0 {
		return c.JSON(http.StatusBadRequest, ErrorMessage("invalid form params"))
	}

	query := `UPDATE "tv" SET 
			"brand"=:brand, "manufacturer"=:manufacturer, "model"=:model, "year"=:year 
			WHERE id=:id`
	_, err = h.DB.NamedExec(query, &tv)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalError())
	}
	return c.JSON(http.StatusCreated, tv)
}

func (h *Handler) DeleteTvById(c echo.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if id < 1 || err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage("Invalid id"))
	}
	query := `DELETE FROM "tv" WHERE id = $1`
	res, err := h.DB.Exec(query, strconv.Itoa(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalError())
	}
	rowsNum, err := res.RowsAffected()
	if rowsNum == 0 && err == nil {
		return c.JSON(http.StatusNotFound, NotFound())
	}
	return
}

func (h *Handler) CreateTv(c echo.Context) (err error) {
	tv := model.Tv{}
	if err = c.Bind(&tv); err != nil {
		return
	}
	tv.Brand.Parse(c.FormValue("brand"))
	if !tv.IsValid() {
		return c.JSON(http.StatusBadRequest, ErrorMessage("invalid form params"))
	}

	query := `INSERT INTO "tv" 
			("id", "brand", "manufacturer", "model", "year") 
			VALUES (DEFAULT, :brand, :manufacturer, :model, :year)`
	res, err := h.DB.NamedExec(query, &tv)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalError())
	}
	id, _ := res.LastInsertId()
	tv.ID = id
	return c.JSON(http.StatusCreated, tv)
}
