'use client';

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { SteganographyService } from '@/service/steganografi';
import type * as React from 'react';
import { useState } from 'react';
import { toast } from 'sonner';

type NullableFile = File | null;

export default function ExtractionSection() {
  const [stegoAudio, setStegoAudio] = useState<NullableFile>(null);
  const [key, setKey] = useState('');
  const [useRandomStart, setUseRandomStart] = useState(false);
  const [outputFileName, setOutputFileName] = useState('');
  const [isProcessing, setIsProcessing] = useState(false);

  const [extractedName, setExtractedName] = useState('secret.bin');
  const [success, setSuccess] = useState<boolean | null>(null);
  const [serverSecretFilename, setServerSecretFilename] = useState<
    string | null
  >(null);
  const [error, setError] = useState<string | null>(null);

  function handleStegoInput(e: React.ChangeEvent<HTMLInputElement>) {
    const f = e.target.files?.[0] || null;
    setStegoAudio(f);
    setSuccess(null);
    setServerSecretFilename(null);
    setError(null);
    setExtractedName('secret.bin');
  }

  const canProcess = !!stegoAudio && !!key && key.length <= 25 && !isProcessing;

  async function onExtract() {
    if (!canProcess || !stegoAudio) return;
    setIsProcessing(true);
    setSuccess(null);
    setServerSecretFilename(null);
    setError(null);
    const pendingOutputName = outputFileName.trim();
    setExtractedName(pendingOutputName || 'secret.bin');

    try {
      const response = await SteganographyService.decode({
        stegoFile: stegoAudio,
        key,
        useRandomStart,
        outputFileName: pendingOutputName || undefined,
      });

      if (!response.success) {
        throw new Error(response.message || 'Failed to decode stego audio');
      }

      setSuccess(true);

      const resolvedServerName =
        response.secretFilename ?? response.secretFileUrl ?? null;
      setServerSecretFilename(resolvedServerName);

      const trimmedOutput = pendingOutputName;

      const fallbackName =
        trimmedOutput ||
        (stegoAudio
          ? `${stegoAudio.name.replace(/\.(mp3|mpeg)$/i, '')}-extracted.bin`
          : 'secret.bin');

      const urlFilename = response.secretFileUrl
        ? response.secretFileUrl.split('/').pop() ?? null
        : null;

      const resolvedName =
        response.secretFilename ?? urlFilename ?? fallbackName;

      setExtractedName(resolvedName);

      toast('Extraction complete', {
        description: response.message,
      });
    } catch (err) {
      const message =
        err instanceof Error ? err.message : 'Failed to decode stego audio';
      setError(message);
      setSuccess(false);
      toast('Extraction failed', {
        description: message,
      });
    } finally {
      setIsProcessing(false);
    }
  }

  function onDownloadSecret() {
    if (!serverSecretFilename) return;

    void (async () => {
      try {
        const blob = await SteganographyService.downloadExtracted(
          serverSecretFilename
        );
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = extractedName || serverSecretFilename;
        document.body.appendChild(a);
        a.click();
        a.remove();
        URL.revokeObjectURL(url);
      } catch (err) {
        const message =
          err instanceof Error ? err.message : 'Failed to download secret file';
        setError(message);
        toast('Download failed', {
          description: message,
        });
      }
    })();
  }

  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Inputs & Uploads</CardTitle>
          <CardDescription>
            Provide the Stego-Audio file (MP3) and the stego key/seed used
            during insertion.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          <div className="grid gap-2">
            <Label htmlFor="stego-audio">Stego Audio File (MP3)</Label>
            <Input
              id="stego-audio"
              type="file"
              accept="audio/mpeg,audio/mp3"
              onChange={handleStegoInput}
            />
            <p className="text-sm text-muted-foreground">
              {stegoAudio
                ? `Selected: ${stegoAudio.name}`
                : 'Choose the MP3 that contains the hidden data.'}
            </p>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="extract-key">Stego Key/Seed (max 25 chars)</Label>
            <Input
              id="extract-key"
              value={key}
              maxLength={25}
              onChange={(e) => setKey(e.target.value)}
              placeholder="Enter your key"
            />
          </div>

          <div className="flex items-center justify-between rounded-lg border border-border/50 bg-muted/30 px-3 py-2">
            <div>
              <Label htmlFor="random-start" className="text-sm">
                Random Start Offset
              </Label>
              <p className="text-xs text-muted-foreground">
                Toggle to randomize the starting sample when decoding.
              </p>
            </div>
            <Switch
              id="random-start"
              checked={useRandomStart}
              onCheckedChange={setUseRandomStart}
            />
          </div>

          <div className="grid gap-2">
            <Label htmlFor="output-name">Desired Output Name (optional)</Label>
            <Input
              id="output-name"
              value={outputFileName}
              onChange={(e) => setOutputFileName(e.target.value)}
              placeholder="secret.bin"
            />
            <p className="text-xs text-muted-foreground">
              Leave empty to let the server choose a filename automatically.
            </p>
          </div>
        </CardContent>
        <CardFooter className="flex items-center justify-end gap-2">
          <Button
            onClick={onExtract}
            disabled={!canProcess}
            aria-busy={isProcessing}
          >
            {isProcessing ? 'Extracting...' : 'Extract Message'}
          </Button>
        </CardFooter>
      </Card>

      {error && (
        <Alert variant="destructive">
          <AlertTitle>Extraction failed</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Results & Download</CardTitle>
          <CardDescription>
            Confirmation and save options for the extracted secret message.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          {success === null && (
            <p className="text-sm text-muted-foreground">
              Extraction results will appear here.
            </p>
          )}
          {success === true && (
            <Alert>
              <AlertTitle>Extraction successful</AlertTitle>
              <AlertDescription>
                The secret message is ready to download.
              </AlertDescription>
            </Alert>
          )}
          {success === false && (
            <Alert variant="destructive">
              <AlertTitle>Extraction failed</AlertTitle>
              <AlertDescription>
                Could not recover the message. Check the key/seed and file.
              </AlertDescription>
            </Alert>
          )}

          <div className="grid gap-2 md:grid-cols-[1fr_auto] md:items-end">
            <div className="grid gap-2">
              <Label htmlFor="secret-filename">Save As</Label>
              <Input
                id="secret-filename"
                value={extractedName}
                onChange={(e) => setExtractedName(e.target.value)}
                placeholder="secret.bin"
              />
            </div>
            <Button onClick={onDownloadSecret} disabled={!serverSecretFilename}>
              Save Secret Message
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
