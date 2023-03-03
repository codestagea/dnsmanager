
import request from "@/utils/http/requests";


export function apiListViews() {
  return request.get(`/v1/view`, {})
}

export function apiUpdateView(id, data) {
  return request.post(`/v1/view/${id}`, { data })

}
export function apiAddView(data) {
  return request.post(`/v1/view`, { data })
}
