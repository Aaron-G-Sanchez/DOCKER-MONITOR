import { Stat } from './stat'

export interface Container {
  id: string
  names: string[]
  state: 'running' | 'exited'
  stats?: Stat
}
