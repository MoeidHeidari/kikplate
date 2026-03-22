import Link from "next/link"
import { HeroSearch } from "@/src/presentation/components/common/HeroSearch"
import { FeaturedPlates } from "@/src/presentation/components/plates/FeaturedPlates"
import { CategoriesGrid } from "@/src/presentation/components/common/CategoriesGrid"
import { StatsBanner } from "@/src/presentation/components/common/StatsBanner"

export default function HomePage() {
  return (
    <div>
      <HeroSearch />
      <FeaturedPlates />
      <CategoriesGrid />
      <StatsBanner />
    </div>
  )
}