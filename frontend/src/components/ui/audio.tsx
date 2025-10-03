import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Slider } from '@/components/ui/slider';
import { Pause, Play, Volume2, VolumeX } from 'lucide-react';
import React, { useEffect, useRef, useState } from 'react';

interface AudioPlayerProps {
  src: string;
}

const AudioPlayer: React.FC<AudioPlayerProps> = ({ src }) => {
  const audioRef = useRef<HTMLAudioElement>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [volume, setVolume] = useState(1); // 0 to 1

  const togglePlayPause = () => {
    const audio = audioRef.current;
    if (audio) {
      if (isPlaying) {
        audio.pause();
      } else {
        audio.play();
      }
      setIsPlaying(!isPlaying);
    }
  };

  const handleSeek = (value: number[]) => {
    const audio = audioRef.current;
    if (audio) {
      audio.currentTime = value[0];
      setCurrentTime(value[0]);
    }
  };

  const handleVolumeChange = (value: number[]) => {
    const newVolume = value[0] / 100;
    const audio = audioRef.current;
    if (audio) {
      audio.volume = newVolume;
      setVolume(newVolume);
    }
  };

  const formatTime = (time: number) => {
    const minutes = Math.floor(time / 60);
    const seconds = Math.floor(time % 60);
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  };

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    const setAudioData = () => {
      setDuration(audio.duration);
    };

    const setTimeUpdate = () => {
      setCurrentTime(audio.currentTime);
    };

    const togglePlay = () => setIsPlaying(true);
    const togglePause = () => setIsPlaying(false);

    audio.addEventListener('loadedmetadata', setAudioData);
    audio.addEventListener('timeupdate', setTimeUpdate);
    audio.addEventListener('play', togglePlay);
    audio.addEventListener('pause', togglePause);

    return () => {
      audio.removeEventListener('loadedmetadata', setAudioData);
      audio.removeEventListener('timeupdate', setTimeUpdate);
      audio.removeEventListener('play', togglePlay);
      audio.removeEventListener('pause', togglePause);
    };
  }, [audioRef]);

  return (
    <Card className="w-full max-w-lg mx-auto border-0">
      <audio ref={audioRef} src={src} preload="metadata" />

      <div className="flex items-center gap-4">
        <span className="text-sm w-10 text-right">
          {formatTime(currentTime)}
        </span>
        <Slider
          value={[currentTime]}
          max={duration || 0}
          step={1}
          onValueChange={handleSeek}
          className="flex-grow"
        />
        <span className="text-sm w-10">{formatTime(duration)}</span>
      </div>

      <div className="flex items-center justify-between">
        <Button
          variant="ghost"
          size="icon"
          onClick={togglePlayPause}
          aria-label={isPlaying ? 'Pause' : 'Play'}
        >
          {isPlaying ? (
            <Pause className="h-6 w-6" />
          ) : (
            <Play className="h-6 w-6" />
          )}
        </Button>

        <div className="flex items-center gap-2 w-1/3">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => handleVolumeChange([volume === 0 ? 100 : 0])}
            aria-label={volume === 0 ? 'Unmute' : 'Mute'}
          >
            {volume === 0 ? (
              <VolumeX className="h-5 w-5" />
            ) : (
              <Volume2 className="h-5 w-5" />
            )}
          </Button>
          <Slider
            value={[volume * 100]}
            max={100}
            step={1}
            onValueChange={handleVolumeChange}
          />
        </div>
      </div>
    </Card>
  );
};

export default AudioPlayer;
