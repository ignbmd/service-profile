package routes

import (
	"github.com/gofiber/fiber/v2"
	"smartbtw.com/services/profile/handlers"
)

func RegisterApiRoute(app *fiber.App) {
	app.Get("/", handlers.HealthCheck)
	app.Get("/students", handlers.GetStudents)
	app.Get("/students/:id", handlers.GetSingleStudent)
	app.Get("/students-report/:id", handlers.GetSingleStudent)
	// app.Get("/students-caching/:id", handlers.GetSingleStudentCaching)
	app.Get("/students-elastic/:id", handlers.GetSingleStudentFromElastic)
	app.Get("/students-with-branch", handlers.GetStudentsByBranchCode)
	app.Get("/students-with-branch/no-limit", handlers.GetStudentsByBranchCodeNoLimit)
	app.Get("/students-with-many-branch", handlers.GetStudentsByArrayBranchCode)
	app.Get("/students-branch", handlers.GetStudentBranch)
	app.Get("/students/joined-class/:id", handlers.GetStudentJoinedClassType)
	app.Get("/branch-list", handlers.GetBranchList)
	app.Post("/student-performa", handlers.GetStudentResultForPerformaSiswa)
	app.Post("/student-performa-uka", handlers.GetStudentResultForPerformaSiswaUKA)

	// Student Report
	app.Get("/student-report/:program/:smId", handlers.FetchRaport)
	app.Get("/student-report-uka/:program/:smId", handlers.FetchRaportUKA)
	app.Get("/student-report-pre-post-test/:program/:smId", handlers.FetchRaportPrePostTest)

	app.Post("/parent-data", handlers.CreateParentData)
	app.Put("/parent-data", handlers.UpdateParentData)

	app.Post("/score-skd", handlers.CreateRecordScore)
	app.Post("/score-skd/year", handlers.GetManyRecordScore)
	app.Get("/score-skd/student/:id", handlers.GetRecordScoreByStudent)
	app.Get("/score-skd/detail/:id", handlers.GetRecordScoreById)
	app.Put("/score-skd/:id", handlers.UpdateRecordScore)
	app.Delete("/score-skd/:id", handlers.DeleteRecordScoreById)

	app.Post("/student-module-progress/", handlers.CreateStudentModuleProgress)
	app.Put("/student-module-progress/:id", handlers.UpdateStudentModuleProgress)
	app.Get("/student-module-progress/detail/:id", handlers.GetStudentModuleProgressBySmartBtwID)
	app.Get("/student-module-progress/task/:task_id", handlers.GetStudentModuleProgressByTaskID)
	app.Get("/student-completed-modules", handlers.GetStudentCompletedModules)

	app.Post("/student-target", handlers.CreateStudentTarget)
	app.Put("/student-target/by-id/:id", handlers.UpdateStudentTargetByID)
	app.Put("/student-target/by-student", handlers.UpdateStudentTargetBySmartbtwID)
	app.Put("/student-target/polbit/by-student", handlers.UpdateStudentPolbitTargetBySmartbtwID)
	app.Put("/student-target/polbit/by-student/bulk", handlers.UpdateBulkStudentPolbitTargetBySmartbtwID)
	app.Put("/student-target/by-student/:smartbtw_id", handlers.UpdateBulkStudentTargetBySmartbtwID)
	app.Get("/student-target/by-student", handlers.GetStudentTarget)
	app.Get("/student-target/by-student-all", handlers.GetAllStudentTarget)
	app.Delete("/student-target/:id", handlers.DeleteStudentTarget)
	app.Get("/student-target/detail/:id", handlers.GetStudentTargetByID)
	app.Get("/student-target/elastic", handlers.GetStudentTargetElastic)
	app.Post("/student-target-cpns", handlers.CreateStudentTargetCPNS)
	app.Put("/student-target/update", handlers.UpdateStudentTarget)
	app.Get("/student-target/competiton-school/:school_origin_id", handlers.GetSchoolCompetition)
	app.Get("/student-target/count-uka-code/:school_origin_id", handlers.CountStudentWithUKACodeBySchoolOriginID)
	app.Post("/student-array/smartbtw-id", handlers.GetStudentBySmartBtwID)

	app.Post("/history-ptk", handlers.CreateHistoryPtk)
	app.Put("/history-ptk/:id", handlers.UpdateHistoryPtk)
	app.Delete("/history-ptk/:id", handlers.DeleteHistoryPtk)
	app.Get("/history-ptk/student/:id", handlers.GetHistoryPtkBySmartBTWID)
	app.Get("/history-ptk/detail/:id", handlers.GetHistoryPtkByID)
	app.Get("/history-ptk/average/:id", handlers.GetStudentAveragePtk)
	app.Get("/history-ptk/last-score/:id", handlers.GetStudentLastScore)
	app.Get("/history-ptk/last-ten-score/:id", handlers.GetLast10StudentScorePtk)
	app.Get("/history-ptk/student-free/:smartbtw_id", handlers.GetStudentFreePTK)
	app.Get("/history-ptk/student-premium/:smartbtw_id", handlers.GetStudentPremiumUKAPTK)
	app.Get("/history-ptk/student-package/:smartbtw_id", handlers.GetStudentPackageUKAPTK)
	app.Post("/history-ptk/all-score", handlers.GetALLStudentScorePtk)
	app.Get("/history-ptk/by-task-id/:task_id", handlers.GetHistoryUKAByTaskID)

	app.Post("/history-cpns", handlers.CreateHistoryCPNS)
	app.Put("/history-cpns/:id", handlers.UpdateHistoryCPNS)
	app.Delete("/history-cpns/:id", handlers.DeleteHistoryCPNS)
	app.Get("/history-cpns/student/:id", handlers.GetHistoryCPNSBySmartBTWID)
	app.Get("/history-cpns/detail/:id", handlers.GetHistoryCPNSByID)
	app.Get("/history-cpns/average/:id", handlers.GetStudentAverageCPNS)
	app.Get("/history-cpns/last-score/:id", handlers.GetStudentLastScoreCPNS)
	app.Get("/history-cpns/last-ten-score/:id", handlers.GetLast10StudentScoreCPNS)
	app.Get("/history-cpns/student-free/:smartbtw_id", handlers.GetStudentFreeCPNS)
	app.Get("/history-cpns/student-premium/:smartbtw_id", handlers.GetStudentPremiumUKACPNS)
	app.Get("/history-cpns/student-package/:smartbtw_id", handlers.GetStudentPackageUKACPNS)
	app.Post("/history-cpns/all-score", handlers.GetALLStudentScoreCPNS)
	app.Get("/history-cpns/by-task-id/:task_id", handlers.GetHistoryCpnsByTaskID)

	app.Post("/history-ptn", handlers.CreateHistoryPtn)
	app.Put("/history-ptn/:id", handlers.UpdateHistoryPtn)
	app.Delete("/history-ptn/:id", handlers.DeleteHistoryPtn)
	app.Get("/history-ptn/detail/:id", handlers.GetHistoryPtnByID)
	app.Get("/history-ptn/student/:id", handlers.GetHistoryPtnBySmartBTWID)
	app.Get("/history-ptn/average/:id", handlers.GetStudentAveragePtn)
	app.Get("/history-ptn/last-score/:id", handlers.GetStudentPtnLastScore)
	app.Get("/history-ptn/last-ten-score/:id", handlers.GetLast10StudentScorePtn)
	app.Get("/history-ptn/student-free/:smartbtw_id", handlers.GetStudentFreePTN)
	app.Get("/history-ptn/student-premium/:smartbtw_id", handlers.GetStudentPremiumUKAPTN)
	app.Get("/history-ptn/student-package/:smartbtw_id", handlers.GetStudentPackageUKAPTN)
	app.Post("/history-ptn/all-score", handlers.GetALLStudentScorePtn)
	app.Get("/history-ptn/by-task-id/:task_id", handlers.GetHistoryPTNByTaskID)

	app.Get("/history-scores/:target_type", handlers.GetHistoryScoreByTargetType)
	app.Get("/history-scores/uka-code-scores/:email", handlers.GetStudentUKAScores)

	app.Get("/wallet/:smartbtw_id/balance", handlers.GetStudentWalletBalance)
	app.Get("/wallet/:smartbtw_id/balance/detail", handlers.GetStudentWalletDetailBalance)
	app.Get("/wallet/:smartbtw_id/history", handlers.GetStudentWalletHistory)
	app.Post("/wallet/:smartbtw_id/charge", handlers.ChargeStudentWallet)
	app.Post("/wallet/check-coin/:smartbtw_id", handlers.CheckCoin)

	app.Post("/avatar", handlers.CreateAvatar)
	app.Get("/avatar/:id", handlers.GetAvatarBySmartbtwID)
	app.Put("/avatar", handlers.UpdateAvatar)
	app.Post("/student-avatar", handlers.GetAvatarBySmartBtwIDAndType)
	app.Delete("/avatar", handlers.DeleteAvatar)

	app.Post("/student-access", handlers.CreateStudentAccess)
	app.Post("/student-access/bulk", handlers.CreateStudentAccessBulk)
	app.Get("/student-access/by-student-id/:id", handlers.GetStudentAccessListBySmartBTWID)
	app.Get("/student-access/by-code/student-id/:id", handlers.GetStudentAccessCode)
	app.Get("/student-access/elastic/student-id/:id", handlers.GetStudentAccessListFromElastic)
	app.Get("/student-access/elastic", handlers.GetStudentAccessListByCodeFromElastic)
	app.Delete("/student-access/student-id/:id", handlers.DeleteStudentAccess)
	app.Delete("/student-access/student-id/bulk/:id", handlers.DeleteStudentAccessBulk)

	app.Get("/class-member/elastic/:id", handlers.GetClassMemberBySmIDElastic)
	app.Get("/class-member/by-classroom-id/:id", handlers.GetClassMemberByClassroomIDElastic)
	app.Get("/classrooms/by-branch/:code", handlers.GetClassroomsByBranchCodes)

	app.Post("/bkn-score", handlers.UpsertBKNScore)
	app.Post("/bkn-score/arr-student-id", handlers.GetBKNScoreByArrayOfSmartbtwID)
	app.Post("/bkn-score/arr-bulk-email", handlers.GetBKNScoreByArrayOfEmail)
	app.Post("/bkn-score/arr-bulk-email/gds", handlers.GetBKNScoreByArrayOfEmailForGDS)
	app.Get("/bkn-score/single/:smartbtw_id/:year", handlers.GetSingleBKNScoreBySMIDAndYear)
	app.Put("/bkn-score/update-survey", handlers.UpdateBKNScoreSurvey)
	app.Put("/bkn-score/update-prodi", handlers.UpdateBKNScoreProdi)

	app.Post("/samapta-score", handlers.UpsertSamaptaScore)
	app.Post("/samapta-score/arr-bulk-email", handlers.GetSamaptaScoreByArrayOfEmail)
	app.Get("/samapta-score/single/:smartbtw_id/:year", handlers.GetSingleSamaptaScoreBySMIDAndYear)

	app.Get("/interview-score/:id", handlers.GetSingleInterviewScoreByID)
	app.Get("/interview-score/student-id/:smartbtw_id/year/:year", handlers.GetSingleInterviewScoreBySMIDAndYear)
	app.Get("/interview-score/session/:session_id/user/:sso_id", handlers.GetInterviewScoresByInterviewSessionIDAndSSOID)
	app.Post("/interview-score", handlers.CreateInterviewScore)
	app.Put("/interview-score/:id", handlers.UpdateInterviewScore)
	app.Post("/interview-score/find-by-emails", handlers.GetInterviewScoreByArrayOfEmail)

	app.Get("/stages/challenge/schools/:program", handlers.GetStagesCompetitionList)
	app.Get("/stages/challenge/schools/:program/ranking/detail/:taskId", handlers.GetStudentSchoolRanking)

	app.Get("/schools/uka/:program/ranking/:taskId", handlers.GetStudentSchoolRankingWithInformation)
	app.Get("/schools/students/:school_id/count", handlers.GetStudentSchoolCount)

	app.Get("/interview-session", handlers.GetAllInterviewSessions)
	app.Get("/interview-session/:id", handlers.GetSingleInterviewSessionByID)
	app.Post("/interview-session", handlers.CreateInterviewSession)
	app.Put("/interview-session/:id", handlers.UpdateInterviewSession)
	app.Delete("/interview-session/:id", handlers.SoftDeleteInterviewSession)

	app.Post("/assessment-screening", handlers.CreateAssessmentScreening)

	app.Get("/progress-result-raport/:program/:type", handlers.GetProgressResultRaport)
	app.Post("/build-raport", handlers.TriggerBuildRaport)
	app.Post("/build-raport-task-id", handlers.TriggerBuildRaportByTaskID)
	app.Post("/get-raport-list", handlers.ListingRaport)
	app.Post("/progress-report", handlers.GetProgressRaport)
	app.Post("/trigger-progress-report", handlers.TriggerBuildProgressResult)
	app.Post("/build-raport-bulk/:program", handlers.BuildRaportBulk)
	app.Post("/build-progress-raport-bulk", handlers.RequestBuildProgressRaportBulk)
	app.Post("/regenerate-raport", handlers.ReGenerateRaportBulk)
}
