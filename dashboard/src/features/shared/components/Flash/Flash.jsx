import React from 'react'
import styles from './Flash.scss'

class Flash extends React.Component {
  componentWillReceiveProps(nextProps) {
    nextProps.messages.forEach((item, key) => {
      if (!item.displayed) {
        this.props.markFlashDisplayed(key)
      }
    })
  }

  render() {
    if (!this.props.messages) {
      return null
    }

    const messages = []
    this.props.messages.forEach((item, key) => {
      messages.push(
        <div className={`${styles.alert} ${styles[item.type]} ${styles.main}`} key={key}>
          <div className={styles.content}>
            {item.title && <div><strong>{item.title}</strong></div>}
            {item.message}
          </div>

          <button type='button' className='close' onClick={() => this.props.dismissFlash(key)}>
            <span>&times;</span>
          </button>
        </div>)
    })

    return (
      <div>
        {messages}
      </div>
    )
  }
}

export default Flash
