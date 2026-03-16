"use client";

import React, { useState, useEffect } from "react";
import { 
  Dialog, 
  DialogContent, 
  DialogHeader, 
  DialogTitle, 
  DialogFooter,
  DialogDescription
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Booking, Slot } from "@/dummy-data/slots";

interface BookingFormProps {
  slot: Slot | null;
  isOpen: boolean;
  onClose: () => void;
  onBook: (booking: Omit<Booking, "id">) => void;
}

export const BookingForm: React.FC<BookingFormProps> = ({ 
  slot, 
  isOpen, 
  onClose, 
  onBook 
}) => {
  const [userName, setUserName] = useState("");
  const [vehicleNumber, setVehicleNumber] = useState("");
  const [duration, setDuration] = useState("60");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!slot) return;
    
    onBook({
      slotId: slot.id,
      userName,
      vehicleNumber,
      startTime: Date.now(),
      duration: parseInt(duration) || 60,
    });
    
    // Reset form
    setUserName("");
    setVehicleNumber("");
    setDuration("60");
    onClose();
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Book Slot {slot?.name}</DialogTitle>
          <DialogDescription>
            Enter your details to reserve this parking spot.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="name" className="text-right">Name</Label>
              <Input
                id="name"
                value={userName}
                onChange={(e) => setUserName(e.target.value)}
                className="col-span-3"
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="vehicle" className="text-right">Vehicle #</Label>
              <Input
                id="vehicle"
                value={vehicleNumber}
                onChange={(e) => setVehicleNumber(e.target.value)}
                placeholder="B 1234 XYZ"
                className="col-span-3"
                required
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="duration" className="text-right">Duration</Label>
              <div className="col-span-3 flex items-center gap-2">
                <Input
                  id="duration"
                  type="number"
                  value={duration}
                  onChange={(e) => setDuration(e.target.value)}
                  min="1"
                  className="flex-1"
                  required
                />
                <span className="text-sm text-zinc-500">min</span>
              </div>
            </div>
          </div>
          <DialogFooter>
            <Button type="submit">Confirm Booking</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};

interface ActiveBookingDetailsProps {
  booking: Booking | null;
  slotName: string | undefined;
  isOpen: boolean;
  onClose: () => void;
  onEndSession: (bookingId: string) => void;
}

export const ActiveBookingDetails: React.FC<ActiveBookingDetailsProps> = ({
  booking,
  slotName,
  isOpen,
  onClose,
  onEndSession
}) => {
  const [timeLeft, setTimeLeft] = useState<{ value: number; isOvertime: boolean }>({ value: 0, isOvertime: false });

  useEffect(() => {
    if (!booking) return;

    const calculateTime = () => {
      const now = Date.now();
      const endTime = booking.startTime + (booking.duration * 60 * 1000);
      const diff = endTime - now;
      
      if (diff > 0) {
        setTimeLeft({ value: Math.floor(diff / 1000), isOvertime: false });
      } else {
        setTimeLeft({ value: Math.floor(Math.abs(diff) / 1000), isOvertime: true });
      }
    };

    calculateTime();
    const interval = setInterval(calculateTime, 1000);
    return () => clearInterval(interval);
  }, [booking]);

  if (!booking) return null;

  const formatTime = (seconds: number) => {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = seconds % 60;
    return `${h > 0 ? `${h}h ` : ""}${m}m ${s}s`;
  };

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Booking Details - {slotName}</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid grid-cols-2 gap-2 text-sm">
            <span className="text-zinc-500 font-medium">User:</span>
            <span>{booking.userName}</span>
            <span className="text-zinc-500 font-medium">Vehicle:</span>
            <span>{booking.vehicleNumber}</span>
            <span className="text-zinc-500 font-medium">Start:</span>
            <span>{new Date(booking.startTime).toLocaleTimeString()}</span>
            <span className="text-zinc-500 font-medium">Duration:</span>
            <span>{booking.duration} minutes</span>
          </div>
          
          <div className={`p-4 rounded-lg flex flex-col items-center justify-center gap-1 ${
            timeLeft.isOvertime ? "bg-red-50 text-red-600 dark:bg-red-900/20" : "bg-emerald-50 text-emerald-600 dark:bg-emerald-900/20"
          }`}>
            <span className="text-xs font-semibold uppercase tracking-wider">
              {timeLeft.isOvertime ? "Overtime" : "Time Remaining"}
            </span>
            <span className="text-3xl font-mono font-bold">
              {formatTime(timeLeft.value)}
            </span>
          </div>
        </div>
        <DialogFooter>
          <Button variant="destructive" onClick={() => onEndSession(booking.id)} className="w-full">
            End Parking Session
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
