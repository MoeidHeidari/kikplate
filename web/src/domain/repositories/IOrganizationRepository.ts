import type {
  Organization,
  CreateOrganizationInput,
  UpdateOrganizationInput,
} from "@/src/domain/entities/Organization"

export interface IOrganizationRepository {
  create(input: CreateOrganizationInput): Promise<Organization>
  listMine(): Promise<Organization[]>
  update(id: string, input: UpdateOrganizationInput): Promise<Organization>
}
