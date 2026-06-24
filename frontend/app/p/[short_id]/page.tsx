import { VoteForm } from "@/components/vote-form";

export default async function PublicPollPage({
  params,
}: {
  params: Promise<{ short_id: string }>;
}) {
  const { short_id } = await params;
  return <VoteForm shortId={short_id} />;
}

export async function generateStaticParams() {
  return [];
}
