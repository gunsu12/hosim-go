// Master Room & Bed Domain Types (internal/master/room)

export interface RoomRecord {
  id: string;
  code: string;
  name: string;
  capacity: number;
  location: string;
  roomType: string;
  unitName: string;
}
