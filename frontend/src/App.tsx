import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

import ExtractionSection from '@/components/stego/extraction-section';
import InsertionSection from '@/components/stego/insertion-section';
import ThemeToggle from '@/components/theme-toggle';

export default function HomePage() {
  return (
    <main className="min-h-dvh">
      <header className="border-b">
        <div className="mx-auto w-full max-w-5xl px-4 py-4 flex items-center justify-between">
          <h1 className="text-xl font-semibold text-balance">
            Audio Steganography Studio
          </h1>
          <div className="flex items-center gap-2">
            <ThemeToggle />
            <Button variant="outline">Help</Button>
          </div>
        </div>
      </header>

      <section className="mx-auto w-full max-w-5xl px-4 py-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-pretty">
              Encode and Decode Secret Messages in Audio
            </CardTitle>
          </CardHeader>
          <CardContent>
            <Tabs defaultValue="insert" className="w-full">
              <TabsList className="grid grid-cols-2 w-full">
                <TabsTrigger value="insert">Insertion (Encode)</TabsTrigger>
                <TabsTrigger value="extract">Extraction (Decode)</TabsTrigger>
              </TabsList>
              <TabsContent value="insert" className="pt-4">
                <InsertionSection />
              </TabsContent>
              <TabsContent value="extract" className="pt-4">
                <ExtractionSection />
              </TabsContent>
            </Tabs>
          </CardContent>
        </Card>
      </section>
    </main>
  );
}
