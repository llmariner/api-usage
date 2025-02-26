import * as fm from "../../fetch.pb";
type Absent<T, K extends keyof T> = {
    [k in Exclude<keyof T, K>]?: undefined;
};
type OneOf<T> = {
    [k in keyof T]?: undefined;
} | (keyof T extends infer K ? (K extends string & keyof T ? {
    [k in K]: T[K];
} & Absent<T, K> : never) : never);
export type UsageRecord = {
    user_id?: string;
    tenant?: string;
    organization?: string;
    project?: string;
    api_key_id?: string;
    api_method?: string;
    status_code?: number;
    timestamp?: string;
    latency_ms?: number;
    details?: UsageDetails;
};
type BaseUsageDetails = {};
export type UsageDetails = BaseUsageDetails & OneOf<{
    create_chat_completion: CreateChatCompletion;
    create_completion: CreateCompletion;
}>;
export type CreateChatCompletion = {
    model_id?: string;
    time_to_first_token_ms?: number;
    prompt_tokens?: number;
    completion_tokens?: number;
};
export type CreateCompletion = {
    model_id?: string;
    time_to_first_token_ms?: number;
    prompt_tokens?: number;
    completion_tokens?: number;
};
export type CreateUsageRequest = {
    records?: UsageRecord[];
};
export type Usage = {};
export declare class CollectionInternalService {
    static CreateUsage(req: CreateUsageRequest, initReq?: fm.InitReq): Promise<Usage>;
}
export {};
