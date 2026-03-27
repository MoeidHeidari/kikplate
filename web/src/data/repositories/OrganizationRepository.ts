import { http } from "./httpClient"
import type { IOrganizationRepository } from "@/src/domain/repositories/IOrganizationRepository"
import type {
  Organization,
  CreateOrganizationInput,
  UpdateOrganizationInput,
} from "@/src/domain/entities/Organization"

class OrganizationRepository implements IOrganizationRepository {
  create(input: CreateOrganizationInput): Promise<Organization> {
    return http.post("/organizations", input)
  }

  listMine(): Promise<Organization[]> {
    return http.get("/organizations/me")
  }

  update(id: string, input: UpdateOrganizationInput): Promise<Organization> {
    return http.put(`/organizations/${id}`, input)
  }
}

export const organizationRepository = new OrganizationRepository()
