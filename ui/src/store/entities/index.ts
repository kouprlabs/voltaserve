import { combineReducers } from 'redux'
import files from './files'
import uploads from './uploads'

export default combineReducers({
  files,
  uploads,
})
