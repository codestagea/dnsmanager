
import request from "@/utils/http/requests";


export function apiListOperationLog(params) {
  return request.get('/v1/operation/logs', { params })
}
