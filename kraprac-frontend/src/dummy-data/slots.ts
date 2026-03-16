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
export const SLOT_COLOR_AVAILABLE = '#10b981';
export const SLOT_COLOR_BOOKED = '#ef4444';
export const SLOT_COLOR_OVERTIME = '#f97316';
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
          { id: 's2', name: 'A2', status: 'available', x: SLOT_WIDTH + SLOT_GAP, y: 0 },
          { id: 's3', name: 'A3', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 2, y: 0 },
          { id: 's4', name: 'A4', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 3, y: 0 },
          { id: 's5', name: 'A5', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 4, y: 0 },
          { id: 's6', name: 'A6', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 5, y: 0 },
          { id: 's7', name: 'A7', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 6, y: 0 },
          { id: 's8', name: 'A8', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 7, y: 0 },
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
          { id: 's11', name: 'B1', status: 'available', x: 0, y: 0 },
          { id: 's12', name: 'B2', status: 'available', x: SLOT_WIDTH + SLOT_GAP, y: 0 },
          { id: 's13', name: 'B3', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 2, y: 0 },
          { id: 's14', name: 'B4', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 3, y: 0 },
          { id: 's15', name: 'B5', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 4, y: 0 },
          { id: 's16', name: 'B6', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 5, y: 0 },
          { id: 's17', name: 'B7', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 6, y: 0 },
          { id: 's18', name: 'B8', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 7, y: 0 },
          { id: 's19', name: 'B9', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 8, y: 0 },
          { id: 's20', name: 'B10', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 9, y: 0 },
        ]
      },
      {
        id: 'r3',
        name: 'Row C',
        x: 50,
        y: 340,
        slots: [
          { id: 's21', name: 'C1', status: 'available', x: 0, y: 0 },
          { id: 's22', name: 'C2', status: 'available', x: SLOT_WIDTH + SLOT_GAP, y: 0 },
          { id: 's23', name: 'C3', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 2, y: 0 },
          { id: 's24', name: 'C4', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 3, y: 0 },
          { id: 's25', name: 'C5', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 4, y: 0 },
          { id: 's26', name: 'C6', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 5, y: 0 },
          { id: 's27', name: 'C7', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 6, y: 0 },
          { id: 's28', name: 'C8', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 7, y: 0 },
          { id: 's29', name: 'C9', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 8, y: 0 },
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
          { id: 's32', name: 'D2', status: 'available', x: SLOT_WIDTH + SLOT_GAP, y: 0 },
          { id: 's33', name: 'D3', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 2, y: 0 },
          { id: 's34', name: 'D4', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 3, y: 0 },
          { id: 's35', name: 'D5', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 4, y: 0 },
          { id: 's36', name: 'D6', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 5, y: 0 },
          { id: 's37', name: 'D7', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 6, y: 0 },
          { id: 's38', name: 'D8', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 7, y: 0 },
          { id: 's39', name: 'D9', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 8, y: 0 },
          { id: 's40', name: 'D10', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 9, y: 0 },
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
          { id: 's43', name: 'E3', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 2, y: 0 },
          { id: 's44', name: 'E4', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 3, y: 0 },
          { id: 's45', name: 'E5', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 4, y: 0 },
          { id: 's46', name: 'E6', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 5, y: 0 },
          { id: 's47', name: 'E7', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 6, y: 0 },
          { id: 's48', name: 'E8', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 7, y: 0 },
          { id: 's49', name: 'E9', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 8, y: 0 },
          { id: 's50', name: 'E10', status: 'available', x: (SLOT_WIDTH + SLOT_GAP) * 9, y: 0 },
        ]
      },
    ]
  },
  {
    id: 'f2',
    name: 'Floor 2',
    blocks: [
      {
        id: 'r6',
        name: 'Row F',
        x: 50,
        y: 60,
        slots: Array.from({ length: 10 }, (_, i) => ({
          id: `s${51 + i}`,
          name: `F${i + 1}`,
          status: 'available',
          x: i * (SLOT_WIDTH + SLOT_GAP),
          y: 0,
        })),
      },
      {
        id: 'r7',
        name: 'Row G',
        x: 50,
        y: 200,
        slots: Array.from({ length: 10 }, (_, i) => ({
          id: `s${61 + i}`,
          name: `G${i + 1}`,
          status: 'available',
          x: i * (SLOT_WIDTH + SLOT_GAP),
          y: 0,
        })),
      },
      {
        id: 'r8',
        name: 'Row H',
        x: 50,
        y: 340,
        slots: Array.from({ length: 10 }, (_, i) => ({
          id: `s${71 + i}`,
          name: `H${i + 1}`,
          status: 'available',
          x: i * (SLOT_WIDTH + SLOT_GAP),
          y: 0,
        })),
      },
      {
        id: 'r9',
        name: 'Row I',
        x: 50,
        y: 480,
        slots: Array.from({ length: 10 }, (_, i) => ({
          id: `s${81 + i}`,
          name: `I${i + 1}`,
          status: 'available',
          x: i * (SLOT_WIDTH + SLOT_GAP),
          y: 0,
        })),
      },
      {
        id: 'r10',
        name: 'Row J',
        x: 50,
        y: 620,
        slots: Array.from({ length: 10 }, (_, i) => ({
          id: `s${91 + i}`,
          name: `J${i + 1}`,
          status: 'available',
          x: i * (SLOT_WIDTH + SLOT_GAP),
          y: 0,
        })),
      },
    ],
  },
  {
    id: 'f3',
    name: 'Floor 3',
    blocks: [
      {
        id: 'r11',
        name: 'Row K',
        x: 50,
        y: 60,
        slots: Array.from({ length: 10 }, (_, i) => ({
          id: `s${101 + i}`,
          name: `K${i + 1}`,
          status: 'available',
          x: i * (SLOT_WIDTH + SLOT_GAP),
          y: 0,
        })),
      },
      {
        id: 'r12',
        name: 'Row L',
        x: 50,
        y: 200,
        slots: Array.from({ length: 10 }, (_, i) => ({
          id: `s${111 + i}`,
          name: `L${i + 1}`,
          status: 'available',
          x: i * (SLOT_WIDTH + SLOT_GAP),
          y: 0,
        })),
      },
      {
        id: 'r13',
        name: 'Row M',
        x: 50,
        y: 340,
        slots: Array.from({ length: 10 }, (_, i) => ({
          id: `s${121 + i}`,
          name: `M${i + 1}`,
          status: 'available',
          x: i * (SLOT_WIDTH + SLOT_GAP),
          y: 0,
        })),
      },
      {
        id: 'r14',
        name: 'Row N',
        x: 50,
        y: 480,
        slots: Array.from({ length: 10 }, (_, i) => ({
          id: `s${131 + i}`,
          name: `N${i + 1}`,
          status: 'available',
          x: i * (SLOT_WIDTH + SLOT_GAP),
          y: 0,
        })),
      },
      {
        id: 'r15',
        name: 'Row O',
        x: 50,
        y: 620,
        slots: Array.from({ length: 10 }, (_, i) => ({
          id: `s${141 + i}`,
          name: `O${i + 1}`,
          status: 'available',
          x: i * (SLOT_WIDTH + SLOT_GAP),
          y: 0,
        })),
      },
    ],
  },
];
