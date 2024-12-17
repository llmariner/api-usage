/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/
import * as fm from "../../fetch.pb";
export class APIUsageService {
    static GetAggregatedSummary(req, initReq) {
        return fm.fetchReq(`/llmariner.apiusage.server.v1.APIUsageService/GetAggregatedSummary`, Object.assign(Object.assign({}, initReq), { method: "POST", body: JSON.stringify(req) }));
    }
    static GetUsageData(req, initReq) {
        return fm.fetchReq(`/llmariner.apiusage.server.v1.APIUsageService/GetUsageData`, Object.assign(Object.assign({}, initReq), { method: "POST", body: JSON.stringify(req) }));
    }
    static ListUsageData(req, initReq) {
        return fm.fetchReq(`/v1/api_usages?${fm.renderURLSearchParams(req, [])}`, Object.assign(Object.assign({}, initReq), { method: "GET" }));
    }
}
