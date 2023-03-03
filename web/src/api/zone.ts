
import request from "@/utils/http/requests";


export function apiListZones(params) {
  return request.get('/v1/zones', { params })
}

export function apiUpdateZone(id, data) {
  return request.post(`/v1/zones/${id}`, { data })

}
export function apiAddZone(data) {
  return request.post(`/v1/zones`, { data })
}
export function apiGetZone(id) {
  return request.get(`/v1/zones/${id}`, {})
}
