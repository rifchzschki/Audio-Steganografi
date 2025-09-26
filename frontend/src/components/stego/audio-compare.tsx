'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { AudioCompareProps } from '@/models/audio';

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
            <audio controls className="w-full">
              <source src={originalUrl} />
              Your browser does not support the audio element.
            </audio>
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
            <audio controls className="w-full">
              <source src={stegoUrl} />
              Your browser does not support the audio element.
            </audio>
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
