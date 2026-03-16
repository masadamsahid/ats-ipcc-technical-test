"use client";

import { useState, useEffect, useMemo } from "react";
import {
  Floor,
  Booking,
  Slot,
  SLOT_COLOR_AVAILABLE,
  SLOT_COLOR_BOOKED
} from "@/dummy-data/slots";
import { getStoredData, saveParkingData } from "@/lib/storage";
import { BookingForm, ActiveBookingDetails } from "@/components/BookingFlow";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  ParkingAreaSquareIcon,
  Car01Icon,
  Calendar03Icon,
  Search01Icon
} from "@hugeicons/core-free-icons";
import { Input } from "@/components/ui/input";
import ParkingMap from "@/components/ParkingMap";
import { HugeiconsIcon } from "@hugeicons/react";

export default function Home() {
  const [data, setData] = useState<{ floors: Floor[]; bookings: Booking[] } | null>(null);
  const [selectedSlot, setSelectedSlot] = useState<Slot | null>(null);
  const [activeBooking, setActiveBooking] = useState<Booking | null>(null);
  const [isBookingOpen, setIsBookingOpen] = useState(false);
  const [isDetailsOpen, setIsDetailsOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");

  // Load data on mount
  useEffect(() => {
    setData(getStoredData());
  }, []);

  // Save data whenever it changes
  useEffect(() => {
    if (data) {
      saveParkingData(data.floors, data.bookings);
    }
  }, [data]);

  const stats = useMemo(() => {
    if (!data) return { total: 0, available: 0, booked: 0 };
    let total = 0;
    let available = 0;
    let booked = 0;

    data.floors.forEach(f => {
      f.blocks.forEach(b => {
        b.slots.forEach(s => {
          total++;
          if (s.status === "available") available++;
          else booked++;
        });
      });
    });

    return { total, available, booked };
  }, [data]);

  const filteredFloors = useMemo(() => {
    if (!data || !searchQuery) return data?.floors;
    const query = searchQuery.toLowerCase();

    return data.floors.map(f => ({
      ...f,
      blocks: f.blocks.map(b => ({
        ...b,
        slots: b.slots.filter(s =>
          s.name.toLowerCase().includes(query) ||
          b.name.toLowerCase().includes(query)
        )
      })).filter(b => b.slots.length > 0)
    })).filter(f => f.blocks.length > 0);
  }, [data, searchQuery]);

  const handleSlotClick = (slot: Slot) => {
    if (slot.status === "available") {
      setSelectedSlot(slot);
      setIsBookingOpen(true);
    } else {
      const booking = data?.bookings.find(b => b.slotId === slot.id);
      if (booking) {
        setActiveBooking(booking);
        setSelectedSlot(slot);
        setIsDetailsOpen(true);
      }
    }
  };
  // ... rest of the functions ...

  const handleBookSlot = (bookingData: Omit<Booking, "id">) => {
    if (!data || !selectedSlot) return;

    const newBooking: Booking = {
      ...bookingData,
      id: `b-${Date.now()}`,
    };

    const newFloors = data.floors.map(f => ({
      ...f,
      blocks: f.blocks.map(b => ({
        ...b,
        slots: b.slots.map(s =>
          s.id === selectedSlot.id ? { ...s, status: "booked" as const } : s
        )
      }))
    }));

    setData({
      floors: newFloors,
      bookings: [...data.bookings, newBooking],
    });

    // Show success (could add toast here)
    console.log("Booking successful!");
  };

  const handleEndSession = (bookingId: string) => {
    if (!data) return;

    const booking = data.bookings.find(b => b.id === bookingId);
    if (!booking) return;

    const newFloors = data.floors.map(f => ({
      ...f,
      blocks: f.blocks.map(b => ({
        ...b,
        slots: b.slots.map(s =>
          s.id === booking.slotId ? { ...s, status: "available" as const } : s
        )
      }))
    }));

    setData({
      floors: newFloors,
      bookings: data.bookings.filter(b => b.id !== bookingId),
    });

    setIsDetailsOpen(false);
    setActiveBooking(null);
  };

  if (!data) return null;

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-black p-4 md:p-8">
      <div className="max-w-6xl mx-auto space-y-8">
        {/* Header */}
        <header className="flex flex-col md:flex-row md:items-end justify-between gap-4">
          <div className="space-y-1">
            <div className="flex items-center gap-2 text-primary">
              <HugeiconsIcon icon={ParkingAreaSquareIcon} size={32} />
              <h1 className="text-3xl font-bold tracking-tight">KrapRac</h1>
            </div>
            <p className="text-zinc-500 dark:text-zinc-400">
              Smart Parking Management System
            </p>
          </div>

          <div className="flex flex-wrap gap-4">
            <Card className="flex items-center px-4 py-2 gap-3 min-w-[120px]">
              <div className="p-2 bg-emerald-100 text-emerald-600 rounded-full dark:bg-emerald-900/30">
                <HugeiconsIcon icon={Car01Icon} size={20} />
              </div>
              <div className="flex flex-col">
                <span className="text-2xl font-bold leading-tight">{stats.available}</span>
                <span className="text-[10px] uppercase font-semibold text-zinc-500">Available</span>
              </div>
            </Card>
            <Card className="flex items-center px-4 py-2 gap-3 min-w-[120px]">
              <div className="p-2 bg-red-100 text-red-600 rounded-full dark:bg-red-900/30">
                <HugeiconsIcon icon={Car01Icon} size={20} />
              </div>
              <div className="flex flex-col">
                <span className="text-2xl font-bold leading-tight">{stats.booked}</span>
                <span className="text-[10px] uppercase font-semibold text-zinc-500">Occupied</span>
              </div>
            </Card>
          </div>
        </header>

        {/* Main Content */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Map Section */}
          <Card className="lg:col-span-2">
            <CardHeader>
              <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                  <CardTitle>Floor Map</CardTitle>
                  <CardDescription>
                    Select an available slot (green) to book or an occupied slot (red) for details.
                  </CardDescription>
                </div>
                <div className="relative w-full md:w-64">
                  <HugeiconsIcon icon={Search01Icon} size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400" />
                  <Input
                    placeholder="Search slot or row..."
                    className="pl-10"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                  />
                </div>
              </div>
            </CardHeader>
            <CardContent>
              {filteredFloors && filteredFloors.length > 0 ? (
                <Tabs defaultValue={filteredFloors[0].id}>
                  <TabsList className="mb-4">
                    {filteredFloors.map(f => (
                      <TabsTrigger key={f.id} value={f.id}>{f.name}</TabsTrigger>
                    ))}
                  </TabsList>
                  {filteredFloors.map(f => (
                    <TabsContent key={f.id} value={f.id}>
                      <ParkingMap floor={f as Floor} onSlotClick={handleSlotClick} />
                    </TabsContent>
                  ))}
                </Tabs>
              ) : (
                <div className="h-[300px] flex flex-col items-center justify-center text-zinc-500 gap-2 border-2 border-dashed rounded-lg">
                  <HugeiconsIcon icon={Search01Icon} size={48} className="opacity-10" />
                  <p>No slots found for "{searchQuery}"</p>
                  <Button variant="link" onClick={() => setSearchQuery("")}>Clear Search</Button>
                </div>
              )}

              <div className="mt-6 flex flex-wrap gap-6 items-center text-sm font-medium border-t pt-4">
                <div className="flex items-center gap-2">
                  <div className="w-4 h-4 rounded-sm" style={{ backgroundColor: SLOT_COLOR_AVAILABLE }}></div>
                  <span>Available</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="w-4 h-4 rounded-sm" style={{ backgroundColor: SLOT_COLOR_BOOKED }}></div>
                  <span>Booked</span>
                </div>
                <div className="text-zinc-400">|</div>
                <div className="text-zinc-500">Total Capacity: {stats.total} Slots</div>
              </div>
            </CardContent>
          </Card>

          {/* Activity Section */}
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Recent Bookings</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {data.bookings.length === 0 ? (
                  <div className="text-center py-8 text-zinc-500 space-y-2">
                    <HugeiconsIcon icon={Calendar03Icon} size={40} className="mx-auto opacity-20" />
                    <p className="text-sm">No active bookings</p>
                  </div>
                ) : (
                  <div className="space-y-3">
                    {data.bookings.slice(-3).reverse().map(b => (
                      <div
                        key={b.id}
                        className="flex items-center justify-between p-3 rounded-lg border bg-zinc-50/50 dark:bg-zinc-900/50 cursor-pointer hover:bg-zinc-100 dark:hover:bg-zinc-900 transition-colors"
                        onClick={() => {
                          setActiveBooking(b);
                          const slot = data.floors.flatMap(f => f.blocks.flatMap(bl => bl.slots)).find(s => s.id === b.slotId);
                          setSelectedSlot(slot || null);
                          setIsDetailsOpen(true);
                        }}
                      >
                        <div className="flex items-center gap-3">
                          <div className="p-2 bg-white dark:bg-zinc-800 rounded border">
                            <HugeiconsIcon icon={Car01Icon} size={16} />
                          </div>
                          <div className="flex flex-col">
                            <span className="font-semibold text-sm">{b.vehicleNumber}</span>
                            <span className="text-[10px] text-zinc-500 uppercase font-bold">Slot {
                              data.floors.flatMap(f => f.blocks.flatMap(bl => bl.slots)).find(s => s.id === b.slotId)?.name
                            }</span>
                          </div>
                        </div>
                        <Badge variant="secondary" className="font-mono text-[10px]">
                          {new Date(b.startTime).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                        </Badge>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>

            <Card className="bg-primary/5 border-primary/20">
              <CardHeader>
                <div className="flex items-center gap-2">
                  <HugeiconsIcon icon={Car01Icon} className="text-primary" size={20} />
                  <CardTitle className="text-lg">Quick Tip</CardTitle>
                </div>
              </CardHeader>
              <CardContent className="text-sm text-zinc-600 dark:text-zinc-400">
                Click on any <span className="text-emerald-600 font-bold">Green</span> slot to make a reservation. For active sessions, click on the <span className="text-red-600 font-bold">Red</span> slot to see details and end the session.
              </CardContent>
            </Card>
          </div>
        </div>
      </div>

      {/* Overlays */}
      <BookingForm
        slot={selectedSlot}
        isOpen={isBookingOpen}
        onClose={() => setIsBookingOpen(false)}
        onBook={handleBookSlot}
      />

      <ActiveBookingDetails
        booking={activeBooking}
        slotName={selectedSlot?.name}
        isOpen={isDetailsOpen}
        onClose={() => setIsDetailsOpen(false)}
        onEndSession={handleEndSession}
      />
    </div>
  );
}
