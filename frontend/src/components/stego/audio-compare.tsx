'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { AudioCompareProps } from '@/models/audio';
import AudioPlayer from '../ui/audio';

export default function AudioCompare({
  originalUrl,
  originalName = 'Cover Audio',
  stegoUrl,
  stegoName = 'Stego Audio',
}: AudioCompareProps) {
  return (
    <div className="grid gap-4 md:grid-cols-2">
      <Card>
        <CardHeader>
          <CardTitle className="text-base">
            {originalName || 'Cover Audio'}
          </CardTitle>
        </CardHeader>
        <CardContent>
          {originalUrl ? (
            <AudioPlayer src={originalUrl} />
          ) : (
            <p className="text-sm text-muted-foreground">
              Upload a cover audio file to preview.
            </p>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">
            {stegoName || 'Stego Audio'}
          </CardTitle>
        </CardHeader>
        <CardContent>
          {stegoUrl ? (
            <AudioPlayer src={stegoUrl} />
          ) : (
            <p className="text-sm text-muted-foreground">
              Stego audio will appear here after insertion.
            </p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
