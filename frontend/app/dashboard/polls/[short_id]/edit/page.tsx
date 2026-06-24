import { EditPollPage } from "@/components/edit-poll-page"

export default async function PollEditPage({
  params,
}: {
  params: Promise<{ short_id: string }>
}) {
  const { short_id } = await params
  return <EditPollPage shortId={short_id} />
}
