import Link from 'next/link';

export default function Redirect({
  url,
  redirectText,
  buttonText,
}: {
  url: string;
  redirectText: string;
  buttonText: string;
}) {
  return (
    <div className="space-x-2">
      <span className="text-xs">{redirectText}</span>
      <Link href={url} className="font-bold text-xs">
        {buttonText}
      </Link>
    </div>
  );
}
