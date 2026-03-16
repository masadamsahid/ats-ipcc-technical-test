"use client";

import React, { useMemo, useState, useEffect } from 'react';
import { Stage, Layer, Rect, Text, Group } from 'react-konva';
import { 
  Floor, 
  Slot, 
  Booking,
  SLOT_COLOR_AVAILABLE, 
  SLOT_COLOR_BOOKED, 
  SLOT_COLOR_OVERTIME,
  SLOT_WIDTH, 
  SLOT_HEIGHT,
  BLOCK_PADDING
} from '@/dummy-data/slots';

interface ParkingMapProps {
  floor: Floor;
  bookings: Booking[];
  onSlotClick: (slot: Slot) => void;
}

const ParkingMap: React.FC<ParkingMapProps> = ({ floor, bookings, onSlotClick }) => {
  const [now, setNow] = useState<number>(0);

  // Update "now" every minute to refresh overtime status
  useEffect(() => {
    // Set initial value after mount to avoid synchronous state update in effect
    const timeout = setTimeout(() => {
      setNow(Date.now());
    }, 0);

    const interval = setInterval(() => {
      setNow(Date.now());
    }, 60000);

    return () => {
      clearTimeout(timeout);
      clearInterval(interval);
    };
  }, []);

  // Calculate stage dimensions based on blocks and slots
  const stageSize = useMemo(() => {
    let maxWidth = 800;
    let maxHeight = 800;

    floor.blocks.forEach(block => {
      block.slots.forEach(slot => {
        const x = block.x + slot.x + SLOT_WIDTH + BLOCK_PADDING;
        const y = block.y + slot.y + SLOT_HEIGHT + BLOCK_PADDING;
        if (x > maxWidth) maxWidth = x;
        if (y > maxHeight) maxHeight = y;
      });
    });

    return { width: maxWidth, height: maxHeight };
  }, [floor]);

  const getSlotColor = (slot: Slot) => {
    if (slot.status === 'available') return SLOT_COLOR_AVAILABLE;
    
    const booking = bookings.find(b => b.slotId === slot.id);
    if (booking) {
      const endTime = booking.startTime + (booking.duration * 60 * 1000);
      if (now > endTime) {
        return SLOT_COLOR_OVERTIME;
      }
    }
    
    return SLOT_COLOR_BOOKED;
  };

  return (
    <div className="overflow-auto border rounded-lg bg-zinc-50 dark:bg-zinc-900/50 p-4">
      <Stage width={stageSize.width} height={stageSize.height}>
        <Layer>
          {floor.blocks.map((block) => (
            <Group key={block.id} x={block.x} y={block.y}>
              <Text
                text={block.name}
                fontSize={16}
                fontStyle="bold"
                fill="#71717a"
                y={-25}
              />
              {block.slots.map((slot) => (
                <Group 
                  key={slot.id} 
                  x={slot.x} 
                  y={slot.y}
                  onClick={() => onSlotClick(slot)}
                  onTap={() => onSlotClick(slot)}
                  style={{ cursor: 'pointer' }}
                >
                  <Rect
                    width={SLOT_WIDTH}
                    height={SLOT_HEIGHT}
                    fill={getSlotColor(slot)}
                    cornerRadius={4}
                    stroke="#00000010"
                    strokeWidth={1}
                  />
                  <Text
                    text={slot.name}
                    width={SLOT_WIDTH}
                    height={SLOT_HEIGHT}
                    align="center"
                    verticalAlign="middle"
                    fill="white"
                    fontSize={12}
                    fontStyle="bold"
                  />
                </Group>
              ))}
            </Group>
          ))}
        </Layer>
      </Stage>
    </div>
  );
};

export default ParkingMap;
