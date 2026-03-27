import type { AppConfig } from "@/src/domain/entities/Config"

export interface IConfigRepository {
  getConfig(): Promise<AppConfig>
}
