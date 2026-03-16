export type SlotStatus = 'available' | 'booked';

export interface Booking {
  id: string; // unique booking id
  slotId: string;
  userName: string;
  vehicleNumber: string;
  startTime: number; // timestamp
  duration: number; // minutes
}

export interface Slot {
  id: string;
  name: string;
  status: SlotStatus;
  x: number;
  y: number;
}

export interface Block {
  id: string;
  name: string;
  x: number;
  y: number;
  slots: Slot[];
}

export interface Floor {
  id: string;
  name: string;
  blocks: Block[];
}

export const STORAGE_KEY = 'kraprac_bookings';
export const SLOT_COLOR_AVAILABLE = '#10b981'; // Emerald 500
export const SLOT_COLOR_BOOKED = '#ef4444';    // Red 500
export const SLOT_WIDTH = 40;
export const SLOT_HEIGHT = 60;
export const SLOT_GAP = 10;
export const BLOCK_PADDING = 20;

export const INITIAL_DATA: Floor[] = [
  {
    id: 'f1',
    name: 'Floor 1',
    blocks: [
      {
        id: 'r1',
        name: 'Row A',
        x: 50,
        y: 60,
        slots: [
          { id: 's1', name: 'A1', status: 'available', x: 0, y: 0 },
          { id: 's2', name: 'A2', status: 'booked', x: SLOT_WIDTH + SLOT_GAP, y: 0 },
          { id: 's3', name: 'A3', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 2, y: 0 },
          { id: 's4', name: 'A4', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 3, y: 0 },
          { id: 's5', name: 'A5', status: 'booked', x: (SLOT_WIDTH + SLOT_GAP) * 4, y: 0 },
          { id: 's6', name: 'A6', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 5, y: 0 },
          { id: 's7', name: 'A7', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 6, y: 0 },
          { id: 's8', name: 'A8', status: 'booked', x: (SLOT_WIDTH + SLOT_GAP) * 7, y: 0 },
          { id: 's9', name: 'A9', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 8, y: 0 },
          { id: 's10', name: 'A10', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 9, y: 0 },
        ]
      },
      {
        id: 'r2',
        name: 'Row B',
        x: 50,
        y: 200,
        slots: [
          { id: 's11', name: 'B1', status: 'booked', x: 0, y: 0 },
          { id: 's12', name: 'B2', status: 'available', x: SLOT_WIDTH + SLOT_GAP, y: 0 },
          { id: 's13', name: 'B3', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 2, y: 0 },
          { id: 's14', name: 'B4', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 3, y: 0 },
          { id: 's15', name: 'B5', status: 'booked', x: (SLOT_WIDTH + SLOT_GAP) * 4, y: 0 },
          { id: 's16', name: 'B6', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 5, y: 0 },
          { id: 's17', name: 'B7', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 6, y: 0 },
          { id: 's18', name: 'B8', status: 'booked', x: (SLOT_WIDTH + SLOT_GAP) * 7, y: 0 },
          { id: 's19', name: 'B9', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 8, y: 0 },
          { id: 's20', name: 'B10', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 9, y: 0 },
        ]
      },
      // add 3 Rows
      {
        id: 'r3',
        name: 'Row C',
        x: 50,
        y: 340,
        slots: [
          { id: 's21', name: 'C1', status: 'available', x: 0, y: 0 },
          { id: 's22', name: 'C2', status: 'available', x: SLOT_WIDTH + SLOT_GAP, y: 0 },
          { id: 's23', name: 'C3', status: 'booked', x: (SLOT_WIDTH + SLOT_GAP) * 2, y: 0 },
          { id: 's24', name: 'C4', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 3, y: 0 },
          { id: 's25', name: 'C5', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 4, y: 0 },
          { id: 's26', name: 'C6', status: 'booked', x: (SLOT_WIDTH + SLOT_GAP) * 5, y: 0 },
          { id: 's27', name: 'C7', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 6, y: 0 },
          { id: 's28', name: 'C8', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 7, y: 0 },
          { id: 's29', name: 'C9', status: 'booked', x: (SLOT_WIDTH + SLOT_GAP) * 8, y: 0 },
          { id: 's30', name: 'C10', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 9, y: 0 },
        ]
      },
      {
        id: 'r4',
        name: 'Row D',
        x: 50,
        y: 480,
        slots: [
          { id: 's31', name: 'D1', status: 'available', x: 0, y: 0 },
          { id: 's32', name: 'D2', status: 'booked', x: SLOT_WIDTH + SLOT_GAP, y: 0 },
          { id: 's33', name: 'D3', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 2, y: 0 },
          { id: 's34', name: 'D4', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 3, y: 0 },
          { id: 's35', name: 'D5', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 4, y: 0 },
          { id: 's36', name: 'D6', status: 'booked', x: (SLOT_WIDTH + SLOT_GAP) * 5, y: 0 },
          { id: 's37', name: 'D7', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 6, y: 0 },
          { id: 's38', name: 'D8', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 7, y: 0 },
          { id: 's39', name: 'D9', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 8, y: 0 },
          { id: 's40', name: 'D10', status: 'booked', x: (SLOT_WIDTH + SLOT_GAP) * 9, y: 0 },
        ]
      },
      {
        id: 'r5',
        name: 'Row E',
        x: 50,
        y: 620,
        slots: [
          { id: 's41', name: 'E1', status: 'available', x: 0, y: 0 },
          { id: 's42', name: 'E2', status: 'available', x: SLOT_WIDTH + SLOT_GAP, y: 0 },
          { id: 's43', name: 'E3', status: 'booked', x: (SLOT_WIDTH + SLOT_GAP) * 2, y: 0 },
          { id: 's44', name: 'E4', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 3, y: 0 },
          { id: 's45', name: 'E5', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 4, y: 0 },
          { id: 's46', name: 'E6', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 5, y: 0 },
          { id: 's47', name: 'E7', status: 'booked', x: (SLOT_WIDTH + SLOT_GAP) * 6, y: 0 },
          { id: 's48', name: 'E8', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 7, y: 0 },
          { id: 's49', name: 'E9', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 8, y: 0 },
          { id: 's50', name: 'E10', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 9, y: 0 },
        ]
      },
    ]
  },
];
