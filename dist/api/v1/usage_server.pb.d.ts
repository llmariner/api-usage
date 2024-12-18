import * as fm from "../../fetch.pb";
export type GetAggregatedSummaryRequest = {
    tenantId?: string;
    startTime?: string;
    endTime?: string;
};
export type Summary = {
    method?: string;
    totalRequests?: string;
    successRequests?: string;
    failureRequests?: string;
    averageLatency?: number;
};
export type AggregatedSummary = {
    summary?: Summary;
    methodSummaries?: Summary[];
};
export type GetUsageDataRequestFilter = {
    apiKeys?: string[];
    models?: string[];
};
export type GetUsageDataRequest = {
    startTime?: string;
    endTime?: string;
    filter?: GetUsageDataRequestFilter;
    after?: string;
    limit?: number;
};
export type UsageData = {
    dataPoints?: UsageDataPoint[];
};
export type UsageDataPoint = {
    userId?: string;
    organization?: string;
    project?: string;
    apiKeyId?: string;
    apiMethod?: string;
    statusCode?: number;
    timestamp?: string;
    latencyMs?: number;
    modelId?: string;
    timeToFirstTokenMs?: number;
    promptTokens?: number;
    completionTokens?: number;
};
export type ListUsageDataRequest = {
    startTime?: string;
    endTime?: string;
};
export type UsageDataByGroup = {
    userId?: string;
    apiKeyId?: string;
    modelId?: string;
    totalRequests?: string;
    avgLatencyMs?: number;
    avgTimeToFirstTokenMs?: number;
    totalPromptTokens?: string;
    totalCompletionTokens?: string;
};
export type ListUsageDataResponse = {
    usages?: UsageDataByGroup[];
};
export declare class APIUsageService {
    static GetAggregatedSummary(req: GetAggregatedSummaryRequest, initReq?: fm.InitReq): Promise<AggregatedSummary>;
    static GetUsageData(req: GetUsageDataRequest, initReq?: fm.InitReq): Promise<UsageData>;
    static ListUsageData(req: ListUsageDataRequest, initReq?: fm.InitReq): Promise<ListUsageDataResponse>;
}
