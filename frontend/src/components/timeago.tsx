import { formatDistanceToNow } from 'date-fns';

interface TimeAgoProps {
  timestamp: string;
}

const TimeAgo: React.FC<TimeAgoProps> = ({ timestamp }) => {
  const timeAgo = formatDistanceToNow(new Date(timestamp), { addSuffix: true });

  return <span>{timeAgo}</span>;
};

export default TimeAgo;
