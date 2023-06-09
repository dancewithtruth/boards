import Link from 'next/link';

export default function Logo({ className }: { className?: string }) {
  return (
    <Link href="/" className={className}>
      Boards
    </Link>
  );
}
