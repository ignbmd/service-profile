package handlers_test

// func Test_UpsertInterviewScore_Success(t *testing.T) {
// 	Init()

// 	payload := request.UpsertInterviewScore{
// 		SmartBtwID: 500,
// 		Penampilan: request.PenampilanScore{
// 			CaraBerpakaian:       3,
// 			CaraDudukDanBerjabat: 3,
// 			PraktekBarisBerbaris: 3,
// 			Total:                9,
// 		},
// 		SikapDanKepribadian: request.SikapDanKepribadianScore{
// 			PenilaianSopanSantun:              3,
// 			KepercayaanDiriDanStabilitasEmosi: 3,
// 			KetahananDiri:                     3,
// 			KelebihanDanKekurangan:            3,
// 			Motivasi:                          3,
// 			Total:                             15,
// 		},
// 		KeluargaDanKemampuanFinansial: request.KeluargaDanKemampuanFinansialScore{
// 			DataKeluargaDanKondisiFinansial: 3,
// 			HubunganDenganTokohNasional:     3,
// 			Total:                           6,
// 		},
// 		SoftSkill: request.SoftSkillScore{
// 			JiwaKepemimpinan:        3,
// 			KemampuanBerkomunikasi:  3,
// 			KemampuanBerbahasaAsing: 3,
// 			Kerjasama:               3,
// 			Total:                   12,
// 		},
// 		HardSkill: request.HardSkillScore{
// 			KemampuanAkademik:      3,
// 			KemampuanMinatDanBakat: 3,
// 			Total:                  3,
// 		},
// 		Total:     100,
// 		CreatedBy: "Interviewer From Unit Test",
// 		Year:      2023,
// 	}

// 	marshalPayload, err := json.Marshal(payload)
// 	assert.Nil(t, err)

// 	app := server.SetupFiber()
// 	request, e := http.NewRequest(
// 		"POST",
// 		"/interview-score",
// 		bytes.NewBuffer(marshalPayload),
// 	)
// 	request.Header.Add("Content-Type", "application/json")
// 	assert.Equal(t, nil, e)

// 	response, err := app.Test(request, -1)
// 	assert.Nil(t, err)

// 	body, err := io.ReadAll(response.Body)
// 	assert.Nil(t, err)

// 	var jBody interface{}
// 	err = sonic.Unmarshal(body, &jBody)
// 	assert.Equal(t, nil, err)
// 	assert.NotNil(t, jBody)

// 	assert.Equal(t, fiber.StatusCreated, response.StatusCode)
// }

// func Test_UpsertInterviewScore_InvalidPayloadDataTypeError(t *testing.T) {
// 	body := []byte(`{
// 	"smartbtw_id": 78987,
// 	"penampilan": {
// 		"cara_berpakaian": 3,
// 		"cara_duduk_dan_berjabat": 3,
// 		"praktek_baris_berbaris": 3,
// 		"total": 9
// 	},
// 	"sikap_dan_kepribadian": {
// 		"penilaian_sopan_santun": 3,
// 		"kepercayaan_diri_dan_stabilitas_emosi": 3,
// 		"ketahanan_diri": 3,
// 		"kelebihan_dan_kekurangan": 3,
// 		"motivasi": 3,
// 		"total": 15
// 	},
// 	"keluarga_dan_kemampuan_finansial": {
// 		"data_keluarga_dan_kondisi_finansial": 3,
// 		"hubungan_dengan_tokoh_nasional": 3,
// 		"total": 6
// 	},
// 	"soft_skill": {
// 		"jiwa_kepemimpinan": 3,
// 		"kemampuan_berkomunikasi": 3,
// 		"kemampuan_berbahasa_asing": 3,
// 		"kerjasama": 3,
// 		"total": 12
// 	},
// 	"hard_skill": {
// 		"kemampuan_akademik": 3,
// 		"kemampuan_minat_dan_bakat": 3,
// 		"total": 6
// 	},
// 	"total": "100",
// 	"created_by": "Another Interviewer",
// 	"year": "2023"
// }`)

// 	app := server.SetupFiber()
// 	request, e := http.NewRequest(
// 		"POST",
// 		"/interview-score",
// 		bytes.NewBuffer(body),
// 	)
// 	request.Header.Add("Content-Type", "application/json")
// 	assert.Equal(t, nil, e)

// 	response, err := app.Test(request, -1)
// 	assert.Nil(t, err)

// 	body, err = io.ReadAll(response.Body)
// 	assert.Nil(t, err)

// 	var jBody interface{}
// 	err = sonic.Unmarshal(body, &jBody)
// 	assert.Equal(t, nil, err)
// 	assert.NotNil(t, jBody)

// 	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
// }

// func Test_UpsertInterviewScore_ValidationError(t *testing.T) {
// 	Init()

// 	payload := request.UpsertInterviewScore{
// 		SmartBtwID: 500,
// 		Penampilan: request.PenampilanScore{
// 			CaraBerpakaian:       3,
// 			CaraDudukDanBerjabat: 3,
// 			PraktekBarisBerbaris: 3,
// 			Total:                9,
// 		},
// 		SikapDanKepribadian: request.SikapDanKepribadianScore{
// 			PenilaianSopanSantun:              3,
// 			KepercayaanDiriDanStabilitasEmosi: 3,
// 			KetahananDiri:                     3,
// 			KelebihanDanKekurangan:            3,
// 			Motivasi:                          3,
// 			Total:                             15,
// 		},
// 		KeluargaDanKemampuanFinansial: request.KeluargaDanKemampuanFinansialScore{
// 			DataKeluargaDanKondisiFinansial: 3,
// 			HubunganDenganTokohNasional:     3,
// 			Total:                           6,
// 		},
// 		SoftSkill: request.SoftSkillScore{
// 			JiwaKepemimpinan:        3,
// 			KemampuanBerkomunikasi:  3,
// 			KemampuanBerbahasaAsing: 3,
// 			Kerjasama:               3,
// 			Total:                   12,
// 		},
// 		HardSkill: request.HardSkillScore{
// 			KemampuanAkademik:      3,
// 			KemampuanMinatDanBakat: 0,
// 			Total:                  3,
// 		},
// 		Year: 2023,
// 	}

// 	marshalPayload, err := json.Marshal(payload)
// 	assert.Nil(t, err)

// 	app := server.SetupFiber()
// 	request, e := http.NewRequest(
// 		"POST",
// 		"/interview-score",
// 		bytes.NewBuffer(marshalPayload),
// 	)
// 	request.Header.Add("Content-Type", "application/json")
// 	assert.Equal(t, nil, e)

// 	response, err := app.Test(request, -1)
// 	assert.Nil(t, err)

// 	body, err := io.ReadAll(response.Body)
// 	assert.Nil(t, err)

// 	var jBody interface{}
// 	err = sonic.Unmarshal(body, &jBody)
// 	assert.Equal(t, nil, err)
// 	assert.NotNil(t, jBody)

// 	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
// }

// func Test_GetInterviewScoreByArrayOfEmail_ValidationError(t *testing.T) {
// 	body := []byte(`{
// 	"year": 2023
// }`)

// 	app := server.SetupFiber()
// 	request, e := http.NewRequest(
// 		"POST",
// 		"/interview-score/find-by-emails",
// 		bytes.NewBuffer(body),
// 	)
// 	request.Header.Add("Content-Type", "application/json")
// 	assert.Equal(t, nil, e)

// 	response, err := app.Test(request, -1)
// 	assert.Nil(t, err)

// 	body, err = io.ReadAll(response.Body)
// 	assert.Nil(t, err)

// 	var jBody interface{}
// 	err = sonic.Unmarshal(body, &jBody)
// 	assert.Equal(t, nil, err)
// 	assert.NotNil(t, jBody)

// 	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
// }

// func Test_GetInterviewScoreByArrayOfEmail_InvalidPayloadDataError(t *testing.T) {
// 	body := []byte(`{
// 	"year": "2023"
// }`)

// 	app := server.SetupFiber()
// 	request, e := http.NewRequest(
// 		"POST",
// 		"/interview-score/find-by-emails",
// 		bytes.NewBuffer(body),
// 	)
// 	request.Header.Add("Content-Type", "application/json")
// 	assert.Equal(t, nil, e)

// 	response, err := app.Test(request, -1)
// 	assert.Nil(t, err)

// 	body, err = io.ReadAll(response.Body)
// 	assert.Nil(t, err)

// 	var jBody interface{}
// 	err = sonic.Unmarshal(body, &jBody)
// 	assert.Equal(t, nil, err)
// 	assert.NotNil(t, jBody)

// 	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
// }

// func Test_GetInterviewScoreByArrayOfEmail_InvalidEmailValidationError(t *testing.T) {
// 	body := []byte(`{
// 	"email": ["contohemailsalah"],
// 	"year": 2023
// }`)

// 	app := server.SetupFiber()
// 	request, e := http.NewRequest(
// 		"POST",
// 		"/interview-score/find-by-emails",
// 		bytes.NewBuffer(body),
// 	)
// 	request.Header.Add("Content-Type", "application/json")
// 	assert.Equal(t, nil, e)

// 	response, err := app.Test(request, -1)
// 	assert.Nil(t, err)

// 	body, err = io.ReadAll(response.Body)
// 	assert.Nil(t, err)

// 	var jBody interface{}
// 	err = sonic.Unmarshal(body, &jBody)
// 	assert.Equal(t, nil, err)
// 	assert.NotNil(t, jBody)

// 	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
// }
