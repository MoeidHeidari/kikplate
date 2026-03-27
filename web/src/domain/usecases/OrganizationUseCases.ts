import type { IOrganizationRepository } from "@/src/domain/repositories/IOrganizationRepository"

export class GetMyOrganizationsUseCase {
  constructor(private readonly repo: IOrganizationRepository) {}
  execute() {
    return this.repo.listMine()
  }
}

export class CreateOrganizationUseCase {
  constructor(private readonly repo: IOrganizationRepository) {}
  execute(input: { name: string; description: string }) {
    return this.repo.create(input)
  }
}

export class UpdateOrganizationUseCase {
  constructor(private readonly repo: IOrganizationRepository) {}
  execute(id: string, input: { name?: string; description?: string }) {
    return this.repo.update(id, input)
  }
}
