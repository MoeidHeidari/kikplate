import { http } from "./httpClient"
import type { IConfigRepository } from "@/src/domain/repositories/IConfigRepository"
import type { AppConfig } from "@/src/domain/entities/Config"

class ConfigRepository implements IConfigRepository {
  getConfig(): Promise<AppConfig> {
    return http.get("/config")
  }
}

export const configRepository = new ConfigRepository()
