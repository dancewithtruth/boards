import avatar from 'gradient-avatar';

const Avatar = ({ id, size = 8 }: { id: string; size?: number }) => {
  const avatarSVG = avatar(id);
  const dataUri = `data:image/svg+xml,${encodeURIComponent(avatarSVG)}`;
  return (
    <div className={`w-${size} h-${size}`}>
      <img className="w-full h-full rounded-full" src={dataUri} alt="Avatar" />
    </div>
  );
};

export default Avatar;
