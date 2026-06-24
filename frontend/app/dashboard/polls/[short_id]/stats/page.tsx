import { PollStats } from "@/components/poll-stats"

export default async function PollStatsPage({
  params,
}: {
  params: Promise<{ short_id: string }>
}) {
  const { short_id } = await params
  return <PollStats shortId={short_id} />
}
