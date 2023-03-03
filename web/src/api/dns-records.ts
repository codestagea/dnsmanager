
import request from "@/utils/http/requests";


export function apiListRecords(zoneId, params) {
  return request.get(`/v1/zone/${zoneId}/records`, { params })
}

export function apiUpdateRecord(zoneId, id, data) {
  return request.post(`/v1/zone/${zoneId}/records/${id}`, { data })

}
export function apiAddRecord(zoneId, data) {
  return request.post(`/v1/zone/${zoneId}/records/`, { data })
}
export function apiGetRecord(zoneId, id) {
  return request.get(`/v1/zone/${zoneId}/records/${id}`, {})
}
