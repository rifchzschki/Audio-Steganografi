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
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Separator } from '@/components/ui/separator';
import { Switch } from '@/components/ui/switch';
import type { LsbBits } from '@/models/stego';
import { SteganographyService } from '@/service/steganografi';
import type * as React from 'react';
import { useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';
import AudioCompare from './audio-compare';

type NullableFile = File | null;

export default function InsertionSection() {
  const [coverFile, setCoverFile] = useState<NullableFile>(null);
  const [secretFile, setSecretFile] = useState<NullableFile>(null);
  const [key, setKey] = useState('');
  const [useEncryption, setUseEncryption] = useState(true);
  const [useRandomStart, setUseRandomStart] = useState(false);
  const [lsbBits, setLsbBits] = useState<LsbBits>(2);
  const [isProcessing, setIsProcessing] = useState(false);

  const [psnr, setPsnr] = useState<number | null>(null);
  const [stegoName, setStegoName] = useState('stego-audio.mp3');
  const [serverStegoFilename, setServerStegoFilename] = useState<string | null>(
    null
  );
  const [stegoStreamUrl, setStegoStreamUrl] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Object URLs
  const coverUrl = useObjectUrl(coverFile);
  const [stegoPreviewUrl, setStegoPreviewUrl] = useState<string | null>(null);

  useEffect(() => {
    setStegoPreviewUrl(stegoStreamUrl);
  }, [stegoStreamUrl]);

  // Capacity (mock-up): Simulate based on audio size and LSB bits to illustrate UI
  const capacityBytes = useMemo(() => {
    if (!coverFile) return 0;
    // NOTE: Real capacity for MP3 would depend on decoded PCM frames & embedding strategy.
    const base = Math.max(coverFile.size * 0.05, 64_000); // ensure non-trivial demo capacity
    return Math.floor(base * lsbBits);
  }, [coverFile, lsbBits]);

  const secretTooLarge = useMemo(() => {
    if (!secretFile) return false;
    return secretFile.size > capacityBytes;
  }, [secretFile, capacityBytes]);

  const canProcess =
    !!coverFile &&
    !!secretFile &&
    !!key &&
    key.length <= 25 &&
    !secretTooLarge &&
    !isProcessing;

  function handleCoverInput(e: React.ChangeEvent<HTMLInputElement>) {
    const f = e.target.files?.[0] || null;
    setCoverFile(f);
    setPsnr(null);
    setServerStegoFilename(null);
    setStegoStreamUrl(null);
    setError(null);
    setStegoName('stego-audio.mp3');
  }

  function handleSecretInput(e: React.ChangeEvent<HTMLInputElement>) {
    const f = e.target.files?.[0] || null;
    setSecretFile(f);
    setPsnr(null);
    setServerStegoFilename(null);
    setStegoStreamUrl(null);
    setError(null);
    setStegoName('stego-audio.mp3');
  }

  async function onInsert() {
    if (!canProcess || !coverFile || !secretFile) return;
    setIsProcessing(true);
    setPsnr(null);
    setServerStegoFilename(null);
    setStegoStreamUrl(null);
    setError(null);

    try {
      const response = await SteganographyService.encode({
        audioFile: coverFile,
        secretFile,
        key,
        lsbBits,
        useEncryption,
        useRandomStart,
      });

      setPsnr(response.psnr ?? null);

      if (response.stegoFileUrl) {
        const fileIdentifier = response.stegoFileUrl;
        const suggestedName = fileIdentifier.split('/').pop() ?? fileIdentifier;
        setServerStegoFilename(fileIdentifier);
        setStegoName(suggestedName);
        const streamUrl =
          SteganographyService.getStegoStreamUrl(fileIdentifier);
        console.log(streamUrl);
        setStegoStreamUrl(streamUrl);
      }

      toast('Insertion complete', {
        description: response.message,
      });
    } catch (err) {
      const message =
        err instanceof Error ? err.message : 'Failed to encode audio';
      setError(message);
      toast('Insertion failed', {
        description: message,
      });
    } finally {
      setIsProcessing(false);
    }
  }

  function onDownloadStego() {
    if (!serverStegoFilename) return;

    void (async () => {
      try {
        const blob = await SteganographyService.downloadStego(
          serverStegoFilename
        );
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = stegoName || serverStegoFilename;
        document.body.appendChild(a);
        a.click();
        a.remove();
        URL.revokeObjectURL(url);
      } catch (err) {
        const message =
          err instanceof Error ? err.message : 'Failed to download stego file';
        setError(message);
        toast('Download failed', { description: message });
      }
    })();
  }

  return (
    <div className="grid gap-6">
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Inputs & Uploads</CardTitle>
          <CardDescription>
            Provide the cover audio (MP3) and the secret message file, then set
            your stego key (max 25 characters).
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          <div className="grid gap-2">
            <Label htmlFor="cover-audio">Cover Audio File (MP3)</Label>
            <Input
              id="cover-audio"
              type="file"
              accept="audio/mpeg,audio/mp3"
              onChange={handleCoverInput}
            />
            <p className="text-sm text-muted-foreground">
              {coverFile
                ? `Selected: ${coverFile.name}`
                : 'Mono or stereo MP3 supported.'}
            </p>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="secret-file">Secret Message File (any type)</Label>
            <Input id="secret-file" type="file" onChange={handleSecretInput} />
            <p className="text-sm text-muted-foreground">
              {secretFile
                ? `Selected: ${secretFile.name}`
                : 'Any file extension is supported.'}
            </p>
          </div>

          <div className="grid gap-2">
            <Label htmlFor="stego-key">Stego Key/Seed (max 25 chars)</Label>
            <Input
              id="stego-key"
              value={key}
              maxLength={25}
              onChange={(e) => setKey(e.target.value)}
              placeholder="Enter your key"
            />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Configuration Options</CardTitle>
          <CardDescription>
            Choose encryption, random start, and number of LSB bits.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          <div className="flex items-center justify-between">
            <div className="space-y-1">
              <Label htmlFor="encryption">
                Use Extended Vigenère Cipher Encryption?
              </Label>
              <p className="text-sm text-muted-foreground">
                Encrypt the secret message before embedding.
              </p>
            </div>
            <Switch
              id="encryption"
              checked={useEncryption}
              onCheckedChange={setUseEncryption}
            />
          </div>

          <Separator />

          <div className="flex items-center justify-between">
            <div className="space-y-1">
              <Label htmlFor="random-start">
                Use Random Start Point for Insertion?
              </Label>
              <p className="text-sm text-muted-foreground">
                Determined from the stego key/seed.
              </p>
            </div>
            <Switch
              id="random-start"
              checked={useRandomStart}
              onCheckedChange={setUseRandomStart}
            />
          </div>

          <Separator />

          <div className="grid gap-2">
            <Label>Multiple-LSB Bits (n-LSB)</Label>
            <RadioGroup
              value={String(lsbBits)}
              onValueChange={(v) => setLsbBits(Number(v) as 1 | 2 | 3 | 4)}
              className="grid grid-cols-4 gap-2"
            >
              <div className="flex items-center space-x-2">
                <RadioGroupItem id="lsb-1" value="1" />
                <Label htmlFor="lsb-1">1 bit</Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem id="lsb-2" value="2" />
                <Label htmlFor="lsb-2">2 bits</Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem id="lsb-3" value="3" />
                <Label htmlFor="lsb-3">3 bits</Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem id="lsb-4" value="4" />
                <Label htmlFor="lsb-4">4 bits</Label>
              </div>
            </RadioGroup>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Capacity Check</CardTitle>
          <CardDescription>
            Evaluate maximum insertion capacity before processing.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-2">
          <div className="text-sm">
            <span className="font-medium">
              Calculated Maximum Insertion Capacity:{' '}
            </span>
            {capacityBytes ? (
              <span>{formatBytes(capacityBytes)}</span>
            ) : (
              <span className="text-muted-foreground">
                Upload a cover audio file to compute capacity.
              </span>
            )}
          </div>
          {secretFile && (
            <div className="text-sm">
              <span className="font-medium">Secret Message Size: </span>
              <span>{formatBytes(secretFile.size)}</span>
            </div>
          )}
          {secretTooLarge && (
            <Alert variant="destructive" className="mt-2">
              <AlertTitle>Secret is too large</AlertTitle>
              <AlertDescription>
                The secret message exceeds the current capacity. Reduce message
                size or increase n-LSB.
              </AlertDescription>
            </Alert>
          )}
        </CardContent>
        <CardFooter className="flex items-center justify-end gap-2">
          <Button
            onClick={onInsert}
            disabled={!canProcess}
            aria-busy={isProcessing}
          >
            {isProcessing ? 'Inserting...' : 'Insert Message'}
          </Button>
        </CardFooter>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Results & Download</CardTitle>
          <CardDescription>
            Review PSNR, compare audio, and save the stego-audio.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
          {error && (
            <Alert variant="destructive">
              <AlertTitle>Insertion failed</AlertTitle>
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}
          <div>
            <Label className="block">PSNR</Label>
            {psnr == null ? (
              <p className="text-sm text-muted-foreground">
                PSNR will appear after insertion.
              </p>
            ) : (
              <div className="text-sm">
                <span className="font-medium">{psnr} dB</span>{' '}
                {psnr < 30 ? (
                  <span className="text-destructive">
                    {' '}
                    (Below 30 dB — significant damage)
                  </span>
                ) : (
                  <span className="text-muted-foreground"> (Acceptable)</span>
                )}
              </div>
            )}
          </div>

          <AudioCompare
            originalUrl={coverUrl || undefined}
            originalName={coverFile?.name}
            stegoUrl={stegoPreviewUrl || undefined}
            stegoName={serverStegoFilename || undefined}
          />

          <div className="grid gap-2 md:grid-cols-[1fr_auto] md:items-end">
            <div className="grid gap-2">
              <Label htmlFor="stego-filename">Save As</Label>
              <Input
                id="stego-filename"
                value={stegoName}
                onChange={(e) => setStegoName(e.target.value)}
                placeholder="stego-audio.mp3"
              />
            </div>
            <Button onClick={onDownloadStego} disabled={!serverStegoFilename}>
              Save Stego-Audio
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

function useObjectUrl(file: File | null) {
  const [url, setUrl] = useState<string | null>(null);
  useEffect(() => {
    if (!file) {
      setUrl(null);
      return;
    }
    const u = URL.createObjectURL(file);
    setUrl(u);
    return () => {
      URL.revokeObjectURL(u);
    };
  }, [file]);
  return url;
}

function formatBytes(bytes: number) {
  if (!bytes) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB'];
  let i = 0;
  let b = bytes;
  while (b >= 1024 && i < units.length - 1) {
    b /= 1024;
    i++;
  }
  return `${b.toFixed(b >= 10 || i === 0 ? 0 : 1)} ${units[i]}`;
}
