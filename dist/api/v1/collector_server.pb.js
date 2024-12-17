/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/
import * as fm from "../../fetch.pb";
export class CollectionInternalService {
    static CreateUsage(req, initReq) {
        return fm.fetchReq(`/llmariner.apiusage.server.v1.CollectionInternalService/CreateUsage`, Object.assign(Object.assign({}, initReq), { method: "POST", body: JSON.stringify(req) }));
    }
}
