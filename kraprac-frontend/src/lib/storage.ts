import { Floor, Booking, INITIAL_DATA, STORAGE_KEY } from "@/dummy-data/slots";

export const getStoredData = (): { floors: Floor[]; bookings: Booking[] } => {
  if (typeof window === "undefined") return { floors: INITIAL_DATA, bookings: [] };
  
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored) {
    try {
      return JSON.parse(stored);
    } catch (e) {
      console.error("Failed to parse stored parking data", e);
    }
  }
  return { floors: INITIAL_DATA, bookings: [] };
};

export const saveParkingData = (floors: Floor[], bookings: Booking[]) => {
  if (typeof window !== "undefined") {
    localStorage.setItem(STORAGE_KEY, JSON.stringify({ floors, bookings }));
  }
};
